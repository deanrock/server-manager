package main

import (
	"../proxy/shell"
	"fmt"
	//"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
	//"errors"
	"log"
	"net"
	"os"
	//"os/user"
	"errors"
	"io/ioutil"
	"encoding/base64"
	"strings"
	//"github.com/kr/pty"
	"golang.org/x/crypto/ssh"
)


func keyAuth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	log.Println(conn.RemoteAddr(), "authenticate user", conn.User(), "with", key.Type())
	fmt.Printf("%s", string(base64.StdEncoding.EncodeToString(key.Marshal())))

	//TODO check user name
	//TODO check key
	//TODO check password for accounts where it's activated

	mykey := "AAAAB3NzaC1yc2EAAAADAQABAAABAQCnhzhmwDZWvRzlZBasnmCh77BB6JpbEQhjv8p02NKFoscVOqB6Ldmb45KB+13QcACkgXfI7ZtxJxgQp6uZjtmrQ3AQc1ukfsQrHVqeLV3h3wZ915L58FvgoAbPFgJj3/JalFWPcc1LQS64wxdphYxprmh5bQ4ZGDVVyPRc7KNXaJ7rWQJTE5rRDaJXmcacda9Ce/FN/8VenQRmYN03Ws9/gyy4fDX9guvrcMyopyOOd41c+J8M4n2byE3R92kf2vHu6dEdDm1eOqg6CU3Pyos+XiBtE1tn15DNYPmzmc/kQlE6jO8mLJWbkoDtRorr/bO6I8Ha/4LR08YvD6Cz+bih"


	if string(base64.StdEncoding.EncodeToString(key.Marshal())) == mykey {
		return nil, nil
	}
	return nil, errors.New("OMG")
}

func handleChannels(chans <-chan ssh.NewChannel) {
	for newChannel := range chans {
		if t := newChannel.ChannelType(); t != "session" {
			newChannel.Reject(ssh.UnknownChannelType, fmt.Sprintf("unknown channel type: %s", t))
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Printf("couldn not accept channel: %s", err)
			continue
		}

		s := shell.Shell{
			LogPrefix: "[ssh]",
		}

		//s.Log("info", "user: %s (%s)", u.Name, u.Uid)

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

		s.AccountName = "flask"
		s.AccountUid = "1005"
		env := "php56"

		ok := false

		fmt.Printf("before req in")

		go func(in <-chan *ssh.Request) {
			for req := range in {
				fmt.Printf("%s", req.Type)
				switch req.Type {
				case "exec":
					ok = true
					s.Cmd =  strings.Split(string(req.Payload[4 : req.Payload[3]+4]), " ")
					s.Tty = false

					s.Environment = strings.Replace(env, "-base-shell", "", 1)
					shell_image, err := s.BuildShellImage(env)

					if err != nil {
						s.LogError(err)
						return
					}

					errs := make(chan error)

					success := make(chan struct{})
					go func() {
						errs <- s.Attach(shell.AttachOptions{
							ShellImage:    shell_image,
							InputStream:   os.Stdin,
							OutputStream:  os.Stdout,
							ErrorStream:   os.Stderr,
							Success:       success,
						})
					}()

					go func() {
						err = <- errs

						if err != nil {
							s.LogError(err)
						}

						channel.Close()
						log.Printf("session (exec, %s (%s)) closed", s.AccountName, s.AccountUid)
					}()
				}
			}
		}(requests)
	}
}

func main() {
	keyPath := "./id_rsa"

	if os.Getenv("KEY_FILE") != "" {
		keyPath = os.Getenv(keyPath)
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

		go handleChannels(chans)
	}
/*
	

	//check if we got "-c" parameter
	if len(os.Args) >= 3 {
		//0: executable name
		//1: "-c"
		//2: command

		if os.Args[1] == "-c" {
			s.Cmd = strings.Split(os.Args[2], " ")
			s.Tty = false
		}
	}

	s.Log("info", "cmd: %s", strings.Join(s.Cmd, " "))

	//environment
	

	if s.Tty {
		//ask user to select environment
		fmt.Printf("Select environment (type the number and press Enter)\n\n")

		for i, image := range(s.ShellImages) {
			fmt.Printf("[%d] %s\n", i+1, image)
		}


		for {
			var i int
			fmt.Printf("Choice: ")
		    _, err = fmt.Scanf("%d", &i)

		    if i >= 1 && i <= len(s.ShellImages) {
		    	env = strings.Replace(
		    		strings.Replace(s.ShellImages[i-1], "-base-shell", "", 1),
		    		"-base-shell", "", 1)
		    	break
		    }
		}
	}

	s.Environment = strings.Replace(env, "-base-shell", "", 1)
	shell_image, err := s.BuildShellImage(env)

	if err != nil {
		s.LogError(err)
		return
	}

	errs := make(chan error)

	success := make(chan struct{})
	go func() {
		errs <- s.Attach(shell.AttachOptions{
			ShellImage:    shell_image,
			InputStream:   os.Stdin,
			OutputStream:  os.Stdout,
			ErrorStream:   os.Stderr,
			Success:       success,
		})
	}()

	if s.Tty {
		{
			file := os.Stdin
			s.InFd = file.Fd()
			term.IsTerminal(s.InFd)
		}

		{
			file := os.Stdout
			s.OutFd = file.Fd()
			term.IsTerminal(s.OutFd)
		}

		oldState, err := term.SetRawTerminal(s.InFd)
		if err != nil {
			return
		}
		defer term.RestoreTerminal(s.InFd, oldState)

		go func() {
			//wait for success to get ContainerID
			<- success

			if err := s.MonitorTtySize(s.ContainerID, false); err != nil {
				s.LogError(fmt.Errorf("Error monitoring TTY size: %s", err))
			}
		}()
	}

	err = <- errs

	if err != nil {
		s.LogError(err)
	}*/
}
