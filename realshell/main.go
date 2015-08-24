package main

import (
	"../proxy/container"
	"errors"
	"fmt"
	"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
	"os"
	"os/user"
	"strings"
)

func main() {
	s := container.Shell{
		LogPrefix: "[shell]",
	}

	u, err := user.Current()

	if err != nil {
		s.LogError(errors.New("cannot get current user"))
		return
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

	//s.Log("info", "cmd: %s", strings.Join(s.Cmd, " "))

	//environment
	s.AccountName = u.Username
	s.AccountUid = u.Uid
	env := "php56"

	if s.Tty {
		//ask user to select environment
		fmt.Printf("Select environment (type the number and press Enter)\n\n")

		for i, image := range s.ShellImages {
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
		errs <- s.Attach(container.AttachOptions{
			ShellImage:   shell_image,
			InputStream:  os.Stdin,
			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
			Success:      success,
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
			<-success

			if err := s.MonitorTtySize(s.ContainerID, false); err != nil {
				s.LogError(fmt.Errorf("Error monitoring TTY size: %s", err))
			}
		}()
	}

	err = <-errs

	if err != nil {
		s.LogError(err)
	}
}
