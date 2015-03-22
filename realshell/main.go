package main

import (
	"../proxy/shell"
	"fmt"
	"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
	"errors"
	"os"
	"os/user"
	"strings"
)

func main() {
	s := shell.Shell{
		LogPrefix: "[shell]",
	}

	u, err := user.Current()

	if err != nil {
		s.LogError(errors.New("cannot get current user"))
		return
	}

	s.Log("info", "user: %s (%s)", u.Name, u.Uid)

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

	s.Log("info", "cmd: %s", strings.Join(s.Cmd, " "))

	//environment
	s.AccountName = u.Username
	s.AccountUid = u.Uid
	env := "nodejs0.12"

	shell_image, err := s.BuildShellImage(env)

	if err != nil {
		s.LogError(err)
		return
	}

	errs := make(chan error)

	go func() {
		errs <- s.Attach(shell.AttachOptions{
			ShellImage:    shell_image,
			InputStream:   os.Stdin,
			OutputStream:  os.Stdout,
			ErrorStream:   os.Stderr,
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

		if err := s.MonitorTtySize(s.ContainerID, false); err != nil {
			s.LogError(fmt.Errorf("Error monitoring TTY size: %s", err))
		}
	}

	err = <- errs

	if err != nil {
		s.LogError(err)
	}
}
