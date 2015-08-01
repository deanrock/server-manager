package main

import (
	"../proxy/shell"
	"../proxy/shared"
	"../proxy/models"
	"fmt"
	"os/exec"
	//"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
	//"errors"
	"log"
	"net"
	//"bytes"
	"strconv"
	"os"
	//"os/user"
	"errors"
	"encoding/binary"
	"io/ioutil"
	"encoding/base64"
	"strings"
	//"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
	"github.com/jinzhu/gorm"
    _ "github.com/mattn/go-sqlite3"
)

var sharedContext *shared.SharedContext

func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	account := models.GetAccountByName(conn.User(), sharedContext)

	log.Println(conn.RemoteAddr(), "authenticate user", conn.User(), "with", key.Type())

	if account == nil {
		err := errors.New(fmt.Sprintf("unknown account %s", conn.User()))
		log.Println(err)
		return nil, err
	}

	key_string := string(base64.StdEncoding.EncodeToString(key.Marshal()))

	keys := models.GetAllUserSSHKeys(sharedContext)

	for _, k := range(keys) {
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
				conn.User(), u.Username, u.Id)

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

		s := shell.Shell{
			LogPrefix: "[ssh]",
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

		s.AccountName = sshConn.User()

		out, err := exec.Command("id","-u",s.AccountName).Output()

		if err != nil {
			return
		}

		uid := strings.Replace(string(out), "\n", "", 1)
		s.AccountUid = uid
		env := "php56"

		w := uint32(0)
		h := uint32(0)

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				fmt.Printf("%s\n", req.Type)
				switch req.Type {
				case "exec":
					ok = true
					s.Cmd =  strings.Split(string(req.Payload[4 : req.Payload[3]+4]), " ")
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
						errs <- s.Attach(shell.AttachOptions{
							ShellImage:    shell_image,
							InputStream:   channel,
							OutputStream:  channel,
							ErrorStream:   channel,
						})
					}()

					go func() {
						err = <- errs
						fmt.Printf("end\n")

						if err != nil {
							s.LogError(err)
						}

						channel.Close()
						log.Printf("session (exec, %s (%s)) closed", s.AccountName, s.AccountUid)
					}()

				case "subsystem":
					ok = true
					subsystem := string(req.Payload[4 : req.Payload[3]+4])
					fmt.Printf("subsystem command %s", subsystem)

					switch subsystem {
					case "sftp":
						//start sftp subsystem
						s.Cmd =  []string{"/usr/lib/openssh/sftp-server"}
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
							errs <- s.Attach(shell.AttachOptions{
								ShellImage:    shell_image,
								InputStream:   channel,
								OutputStream:  channel,
								ErrorStream:   channel,
							})
						}()

						go func() {
							err = <- errs
							fmt.Printf("end\n")

							if err != nil {
								s.LogError(err)
							}

							channel.Close()
							log.Printf("session (subsystem %s, %s (%s)) closed", subsystem, s.AccountName, s.AccountUid)
						}()
					default:
						log.Printf("session (subsystem, %s (%s)) - wrong command %s", s.AccountName, s.AccountUid, subsystem)
						channel.Close()
					}

				case "shell":
					ok = true
					s.Cmd = []string {"/bin/bash"} //  strings.Split(string(req.Payload[4 : req.Payload[3]+4]), " ")
					s.Tty = true

					printChoiceMenu(channel, s)

					for {
						var i int
						channel.Write([]byte("\r\nChoice: "))
						
						selected := make([]byte, 10)
						n, err := channel.Read(selected)

						if err != nil {
							log.Printf("error reading from channel: %s", err)
							return
						}

						fmt.Printf("%s\n", selected)
						//parse number
						d, err := strconv.ParseInt(string(selected[0:n]), 0, 64)

						if err != nil {
							continue
						}

						i = int(d)

					    if i >= 1 && i <= len(s.ShellImages) {
					    	env = strings.Replace(
					    		strings.Replace(s.ShellImages[i-1], "-base-shell", "", 1),
					    		"-base-shell", "", 1)

					    	channel.Write([]byte("\r\n"))
					    	break
					    }
					}

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
						errs <- s.Attach(shell.AttachOptions{
							ShellImage:    shell_image,
							InputStream:   channel,
							OutputStream:  channel,
							ErrorStream:   channel,
							Success:       success,
							Detach:        detach,
						})
					}()

					go func() {
						<- success
						resizeTty(s, h, w)
					}()

					go func() {
						err = <- errs
						fmt.Printf("end\n")

						if err != nil {
							s.LogError(err)
						}

						channel.Close()
						log.Printf("session (shell, %s (%s)) closed", s.AccountName, s.AccountUid)
					}()
					fmt.Printf("yo\n")
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

func printChoiceMenu(channel ssh.Channel, s shell.Shell) {
	//ask user to select environment
	//channel.Write([]byte{0x1b,'[','H',0x1b,'[','J'})
	channel.Write([]byte("Select environment (type the number and press Enter)\r\n\r\n"))

	for i, image := range(s.ShellImages) {
		channel.Write([]byte(fmt.Sprintf("[%d] %s\r\n", i+1, image)))
	}

	//channel.Write([]byte("Choice: "))
}

func resizeTty(shell shell.Shell, height uint32, width uint32) {
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

func main() {
	f, err := os.OpenFile("/var/log/manager/ssh-server.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	sharedContext = &shared.SharedContext{}

    //sqlite
    db, err := gorm.Open("sqlite3", "../manager/db.sqlite3")

    if err != nil {
        log.Fatal("database error", err)
    }

    sharedContext.PersistentDB = db

	keyPath := "./id_rsa"

	if os.Getenv("KEY_FILE") != "" {
		keyPath = os.Getenv("KEY_FILE")
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
	}
	config.AddHostKey(keySigner)

	port := "2222"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	socket, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := socket.Accept()
		if err != nil {
			log.Printf("failed to accept the connection: %s", err)
			continue
		}

		sshConn, chans, _, err := ssh.NewServerConn(conn, &config)
		if err != nil {
			log.Printf("handshake failed: %s", err)
			continue
		}

		log.Printf("new ssh connection from %s (%s)", sshConn.RemoteAddr(), sshConn.ClientVersion())

		go handleChannels(sshConn, chans)
	}
}
