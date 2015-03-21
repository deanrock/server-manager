package main

import (
	"../proxy/shell"
	"fmt"
	"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

func main() {
	u, e := user.Current()

	if e != nil {
		fmt.Printf("error")
		return
	}

	s := shell.Shell{}
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

	//go-dockerclient
	endpoint := "unix:///var/run/docker.sock"
	s.DockerClient, e = docker.NewClient(endpoint)

	if e != nil {
		return
	}

	//environment
	account := u.Username
	env := "php56-base-shell"

	found := false
	for _, e := range s.ShellImages {
		if e == env {
			found = true
		}
	}

	if !found {
		fmt.Printf("not found")
		return
	}

	s.Mylog(strings.Join(s.Cmd, " "))

	out, err := exec.Command("id", "-u", account).Output()

	if err != nil {
		return
	}

	container, err := s.DockerClient.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			OpenStdin: true,
			Tty:       s.Tty,
			Cmd:       s.Cmd,
			Image:     "manager/" + env,
			Env:       []string{"USER=" + account, "USERID=" + strings.Replace(string(out), "\n", "", 1)},
		},
	})

	defer func() {
		return
		//cleanup
		log.Printf("cleanup shell %s %s", account, env)
		err := s.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})

		if err != nil {
			log.Printf("error while cleaning up ", err)
		}
	}()

	if err != nil {
		log.Printf("cannot create container ", err)
		return
	}

	err = s.DockerClient.StartContainer(container.ID,
		&docker.HostConfig{
			Binds:      []string{"/home/" + account + ":/home/" + account},
			ExtraHosts: []string{"mysql:172.17.42.1"},
		})

	if err != nil {
		log.Printf("cannot start container ", err)
		return
	}

	//r, w := io.Pipe()
	//stdinR, stdinW := io.Pipe()

	//os.Stdout =
	var errs error
	go func() {
		errs = s.DockerClient.AttachToContainer(docker.AttachToContainerOptions{
			Container:    container.ID,
			OutputStream: os.Stdout,
			ErrorStream:  os.Stderr,
			InputStream:  os.Stdin,
			Stdout:       true,
			Stdin:        true,
			Stderr:       true,
			Stream:       true,
			RawTerminal:  s.Tty,
		})
	}()
	if err != nil {
		log.Printf("cannot attach to container ", err)
		return
	}

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

		if err := s.MonitorTtySize(container.ID, false); err != nil {
			log.Printf("Error monitoring TTY size: %s", err)
		}
	}

	/*if attchErr := <-cErr; attchErr != nil {
		return attchErr
	}
	_, status, err := getExitCode(cli, cmd.Arg(0))
	if err != nil {
		return err
	}
	if status != 0 {
		return &utils.StatusError{StatusCode: status}
	}*/

	/*func() {
	_, err = io.Copy(os.Stdin, )
	}()

	func() {
	_, err = io.Copy(stdout, br)
	}()*/

	//os.Stdout = w

	/*go func(reader io.Reader) {
		/*for {
			data := make([]byte, 100)
			n, err := reader.Read(data)

			if err != nil {
				log.Printf("error while reading from docker stream ", err)
			}



		}
	}(r)*/

	for {

		time.Sleep(1 * time.Minute)
		//stdinW.Write([]byte{0x44})
	}
}
