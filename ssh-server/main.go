package ssh_server

import (
	"fmt"
	"os/exec"

	"../proxy/container"
	"../proxy/models"
	"../proxy/shared"
	"github.com/fsouza/go-dockerclient"
	//"errors"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/ssh"
)

var sharedContext *shared.SharedContext

func getUserAndEnvironment(username string) (string, string) {
	s := strings.Split(username, "+")

	if len(s) > 1 {
		return s[0], s[1]
	}

	return s[0], ""
}

func passwordAuth(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	u, _ := getUserAndEnvironment(conn.User())
	account := models.GetAccountByName(u, sharedContext)

	log.Println(conn.RemoteAddr(), "authenticate user", u, "with password")

	if account == nil {
		err := errors.New(fmt.Sprintf("unknown account %s", u))
		log.Println(err)
		return nil, err
	}

	var passwords []models.SSHPassword
	sharedContext.PersistentDB.Where("account_id = ?", account.Id).Find(&passwords)

	for _, p := range passwords {
		if p.Password == string(password) && len(string(password)) >= 10 {
			return nil, nil
		}
	}

	return nil, errors.New("cannot authenticate user")
}

func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	user, _ := getUserAndEnvironment(conn.User())

	account := models.GetAccountByName(user, sharedContext)

	log.Println(conn.RemoteAddr(), "authenticate user", user, "with", key.Type())

	if account == nil {
		err := errors.New(fmt.Sprintf("unknown account %s", user))
		log.Println(err)
		return nil, err
	}

	key_string := string(base64.StdEncoding.EncodeToString(key.Marshal()))

	keys := models.GetAllUserSSHKeys(sharedContext)

	for _, k := range keys {
		mykey := k.SSHKey
		s := strings.Split(mykey, " ")

		if len(s) > 1 && s[0] == "ssh-rsa" {
			mykey = s[1]
		}

		if key_string == mykey {
			u, err := models.FindUserById(sharedContext, k.User_id)

			if err != nil {
				return nil, err
			}

			log.Printf("authenticated %s by ssh key from %s (user id: %d)",
				user, u.Username, u.Id)

			var userAccess models.UserAccess
			if err := sharedContext.PersistentDB.Where("user_id = ? AND account_id = ?", u.Id, account.Id).First(&userAccess).Error; err != nil {
				log.Printf("user %s (%d) doesn't have access to this account %s", u.Username, u.Id, account.Name)
				return nil, errors.New("no access")
			}

			if !userAccess.SshAccess {
				log.Printf("user %s (%d) doesn't have SSH access to this account %s", u.Username, u.Id, account.Name)
				return nil, errors.New("no SSH access")
			}

			return nil, nil //success
		}
	}

	//TODO check password for accounts where it's activated

	return nil, errors.New("cannot authenticate user")
}

func handleChannels(sshConn *ssh.ServerConn, chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		if t := newChannel.ChannelType(); t != "session" {
			newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("could not accept channel: %s", err)
			continue
		}

		s := container.Shell{
			LogPrefix:     "[ssh]",
			SharedContext: sharedContext,
		}

		//go-dockerclient
		endpoint := "unix:///var/run/docker.sock"
		s.DockerClient, err = docker.NewClient(endpoint)

		if err != nil {
			s.LogError(errors.New("cannot connect to docker client"))
			return
		}

		s.GetDockerImages()

		s.Cmd = []string{"/bin/bash"}
		s.Tty = true

		u, e := getUserAndEnvironment(sshConn.User())
		s.AccountName = u

		out, err := exec.Command("id", "-u", s.AccountName).Output()

		if err != nil {
			returnExitCode(255, channel)
			channel.Close()
			log.Println("user doesn't exist")
			continue
		}

		uid := strings.Replace(string(out), "\n", "", 1)
		s.AccountUid = uid

		// determine environment
		env := "php56"

		if e != "" {
			env = e
		}

		log.Println("wanted environment:", env)

		if s.ImageAllowed(env) != nil {
			returnExitCode(255, channel)
			channel.Close()
			log.Println("wrong environemnt")
			continue
		}

		w := uint32(0)
		h := uint32(0)

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				log.Printf("%s\n", req.Type)
				switch req.Type {
				case "exec":
					ok = true
					s.Cmd = strings.Split(string(req.Payload[4:req.Payload[3]+4]), " ")
					s.Tty = false

					s.Environment = strings.Replace(env, "-base-shell", "", 1)
					shell_image, err := s.BuildShellImage(env)

					log.Printf("command %s", strings.Join(s.Cmd, " "))

					if err != nil {
						s.LogError(err)
						return
					}

					errs := make(chan error)

					r, w := io.Pipe()

					go func() {
						errs <- s.Attach(container.AttachOptions{
							ShellImage:   shell_image,
							InputStream:  r,
							OutputStream: channel,
							ErrorStream:  channel,
						})
					}()

					go io.Copy(w, channel)

					go func() {
						err = <-errs
						log.Printf("end\n")

						if err != nil {
							s.LogError(err)
						}

						returnExitCode(s.ExitCode, channel)
						channel.Close()
						log.Printf("session (exec, %s (%s)) closed", s.AccountName, s.AccountUid)
					}()

				case "subsystem":
					ok = true
					subsystem := string(req.Payload[4 : req.Payload[3]+4])
					log.Printf("subsystem command %s", subsystem)

					switch subsystem {
					case "sftp":
						//start sftp subsystem
						s.Cmd = []string{"/usr/lib/openssh/sftp-server"}
						s.Tty = false

						s.Environment = strings.Replace(env, "-base-shell", "", 1)
						shell_image, err := s.BuildShellImage(env)

						log.Printf("command %s", strings.Join(s.Cmd, " "))

						if err != nil {
							s.LogError(err)
							return
						}

						errs := make(chan error)

						go func() {
							errs <- s.Attach(container.AttachOptions{
								ShellImage:   shell_image,
								InputStream:  channel,
								OutputStream: channel,
								ErrorStream:  channel,
							})
						}()

						go func() {
							err = <-errs
							log.Printf("end\n")

							if err != nil {
								s.LogError(err)
							}

							returnExitCode(s.ExitCode, channel)
							channel.Close()
							log.Printf("session (subsystem %s, %s (%s)) closed", subsystem, s.AccountName, s.AccountUid)
						}()
					default:
						log.Printf("session (subsystem, %s (%s)) - wrong command %s", s.AccountName, s.AccountUid, subsystem)
						returnExitCode(127, channel)
						channel.Close()
					}

				case "shell":
					ok = true
					s.Cmd = []string{"/bin/bash"} //  strings.Split(string(req.Payload[4 : req.Payload[3]+4]), " ")
					s.Tty = true

					s.Environment = strings.Replace(env, "-base-shell", "", 1)
					shell_image, err := s.BuildShellImage(env)

					if err != nil {
						s.LogError(err)
						return
					}

					errs := make(chan error)
					detach := make(chan error)

					success := make(chan struct{})
					go func() {
						errs <- s.Attach(container.AttachOptions{
							ShellImage:   shell_image,
							InputStream:  channel,
							OutputStream: channel,
							ErrorStream:  channel,
							Success:      success,
							Detach:       detach,
						})
					}()

					go func() {
						<-success
						resizeTty(s, h, w)
					}()

					go func() {
						err = <-errs
						log.Printf("end\n")

						if err != nil {
							s.LogError(err)
						}

						returnExitCode(s.ExitCode, channel)
						channel.Close()
						log.Printf("session (shell, %s (%s)) closed", s.AccountName, s.AccountUid)
					}()
				case "pty-req":
					// Responding 'ok' here will let the client
					// know we have a pty ready for input
					ok = true
					// Parse body...
					termLen := req.Payload[3]
					termEnv := string(req.Payload[4 : termLen+4])
					w, h = parseDims(req.Payload[termLen+4:])
					resizeTty(s, h, w)
					log.Printf("pty-req '%s'", termEnv)
				case "window-change":
					w, h := parseDims(req.Payload)
					resizeTty(s, h, w)
					continue //no response
				}

				if !ok {
					log.Printf("declining %s request...", req.Type)
				}

				req.Reply(ok, nil)
			}
		}(requests)
	}
}

func resizeTty(shell container.Shell, height uint32, width uint32) {
	if shell.ContainerID == "" {
		log.Printf("cannot resize tty; containerid empty")
		return
	}

	shell.ResizeTtyTo(shell.ContainerID, int(height), int(width))
}

// parseDims extracts two uint32s from the provided buffer.
func parseDims(b []byte) (uint32, uint32) {
	w := binary.BigEndian.Uint32(b)
	h := binary.BigEndian.Uint32(b[4:])
	return w, h
}

func returnExitCode(code int, channel ssh.Channel) {
	log.Println("exit code", code)

	// return exit-status for ssh client
	// code can be forced to uint32, because bash specifies 8-bit uint as exit-code type
	bs := make([]byte, 4)
	binary.BigEndian.PutUint32(bs, uint32(code))
	channel.SendRequest("exit-status", false, bs)
}

func handleConnection(conn net.Conn, config ssh.ServerConfig) {
	sshConn, chans, _, err := ssh.NewServerConn(conn, &config)
	if err != nil {
		log.Printf("handshake failed: %s", err)
		return
	}

	log.Printf("new ssh connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

	handleChannels(sshConn, chans)
}

func Start() {
	f, err := os.OpenFile("/var/log/manager/ssh-server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	sharedContext = &shared.SharedContext{}

	//sqlite
	db, err := gorm.Open("sqlite3", "db.sqlite3")

	if err != nil {
		log.Fatal("database error", err)
	}

	sharedContext.PersistentDB = db

	err = models.ParseImages(sharedContext)
	if err != nil {
		log.Fatalf("error encountered while parsing images: %s", err)
	}

	sharedContext.InitConfig("config.json")

	keyPath := "./id_rsa"

	if sharedContext.Config.SSHConfig.KeyPath != "" {
		keyPath = sharedContext.Config.SSHConfig.KeyPath
	}

	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(err)
	}

	keySigner, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		panic(err)
	}

	config := ssh.ServerConfig{
		PublicKeyCallback: keyAuth,
		PasswordCallback:  passwordAuth,
		Config: ssh.Config{
			Ciphers: []string{
				"aes128-ctr", "aes192-ctr", "aes256-ctr",
				"aes128-gcm@openssh.com",
				"arcfour256", "arcfour128",
				"aes128-cbc", // this insecure crypto is enabled because Aptana Studio doesn't
				// support anything more secure :/
			},
		},
	}
	config.AddHostKey(keySigner)

	port := "2222"
	if sharedContext.Config.SSHConfig.Port != 0 {
		port = fmt.Sprintf("%d", sharedContext.Config.SSHConfig.Port)
	}

	socket, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("failed to accept the connection: %s", err)
			continue
		}

		log.Printf("accept from: %s", conn.RemoteAddr())

		go handleConnection(conn, config)
	}
}
