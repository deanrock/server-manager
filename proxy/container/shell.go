package container

import (
	"../models"
	"../shared"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/docker/docker/pkg/signal"
	"github.com/docker/docker/pkg/term"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	gosignal "os/signal"
	"path"
	"strings"
	"time"
)

type Shell struct {
	SharedContext *shared.SharedContext

	LogPrefix    string
	AccountName  string
	AccountUid   string
	Images       []string
	ShellImages  []string
	Cmd          []string
	Tty          bool
	DockerClient *docker.Client
	ContainerID  string
	Environment  string

	// inFd holds file descriptor of the client's STDIN, if it's a valid file
	InFd uintptr
	// outFd holds file descriptor of the client's STDOUT, if it's a valid file
	OutFd uintptr
}

func (s *Shell) BuildShellImage(env string) (string, error) {
	base_image := fmt.Sprintf("%s-base-shell", env)

	found := false
	for _, e := range s.ShellImages {
		if e == base_image {
			found = true
		}
	}

	if !found {
		return "", fmt.Errorf("base shell %s not found", base_image)
	}

	//build shell image if it doesnt exist
	shell_image := fmt.Sprintf("shell-%s-%s", s.AccountName, env)

	for _, e := range s.Images {
		if e == shell_image {
			return shell_image, nil
		}
	}

	//couldn't find image, build it
	temp, err := ioutil.TempDir("", "manager-")

	if err != nil {
		return "", fmt.Errorf("cannot create temp folder: %s", err)
	}

	defer os.RemoveAll(temp)

	image_folder := fmt.Sprintf("/home/manager/server-manager/images/%s/",
		base_image)

	if _, err := os.Stat(image_folder); os.IsNotExist(err) {
		return "", fmt.Errorf("no such file or directory: %s", image_folder)
	}

	files, _ := ioutil.ReadDir(image_folder)
	for _, f := range files {
		//copy file
		out_file, err := os.Create(path.Join(temp, f.Name()))

		if err != nil {
			return "", fmt.Errorf("cannot create temp file %s", err)
		}

		defer out_file.Close()

		in_file, err := os.Open(path.Join(image_folder, f.Name()))
		if err != nil {
			return "", fmt.Errorf("cannot open file to read %s", err)
		}
		defer in_file.Close()

		scanner := bufio.NewScanner(in_file)
		for scanner.Scan() {
			_, err := io.WriteString(out_file, scanner.Text()+"\n")
			if err != nil {
				return "", fmt.Errorf("error writing temp file %s", err)
			}
		}

		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("error scanning file %s", err)
		}

		if f.Name() == "Dockerfile" {
			_, err := io.WriteString(out_file,
				fmt.Sprintf("RUN echo \"%s:x:%s:\" >> /etc/group && echo \"%s:x:%s:%s:,,,:/home/%s:/bin/bash\" >> /etc/passwd\n\nUSER %s\n",
					s.AccountName,
					s.AccountUid,
					s.AccountName,
					s.AccountUid,
					s.AccountUid,
					s.AccountName,
					s.AccountName))

			if err != nil {
				return "", fmt.Errorf("error writing to Dockerfile %s", err)
			}
		}
	}

	var buf bytes.Buffer
	err = s.DockerClient.BuildImage(docker.BuildImageOptions{
		Name:         fmt.Sprintf("manager/%s", shell_image),
		ContextDir:   temp,
		OutputStream: &buf,
	})

	if err != nil {
		return "", fmt.Errorf("building image error %s", err)
	}

	return shell_image, nil
}

type AttachOptions struct {
	ShellImage   string
	OutputStream io.Writer
	ErrorStream  io.Writer
	InputStream  io.Reader
	Success      chan struct{}
	Detach       chan error
}

func (shell *Shell) CreateContainer(shellImage string) (*docker.Container, error) {
	container, err := shell.DockerClient.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			OpenStdin: true,
			Tty:       shell.Tty,
			Cmd:       shell.Cmd,
			Image:     "manager/" + shellImage,
			Hostname:  shell.Environment,
		},
	})

	if err == nil {
		shell.ContainerID = container.ID
	}

	return container, err
}

func (shell *Shell) StartContainer() error {
	if shell.SharedContext != nil {
		account := models.GetAccountByName(shell.AccountName, shell.SharedContext)
		shell.GetDockerImages()

		err := StartContainer(account, shell.SharedContext, shell.DockerClient, shell.ContainerID)
		return err
	}

	return errors.New("no shared context provided")
}

func (shell *Shell) RemoveContainer() error {
	if shell.ContainerID != "" {
		err := shell.DockerClient.RemoveContainer(docker.RemoveContainerOptions{
			ID:    shell.ContainerID,
			Force: true,
		})

		return err
	}

	return nil
}

func (shell *Shell) Attach(options AttachOptions) error {
	container, err := shell.CreateContainer(options.ShellImage)
	if err != nil {
		return fmt.Errorf("couldn't create container %s (image: %s)", err, options.ShellImage)
	}

	defer func() {
		shell.Log("info", "cleanup shell %s %s", shell.AccountName, options.ShellImage)

		if container != nil {
			err := shell.RemoveContainer()

			if err != nil {
				shell.LogError(fmt.Errorf("error while cleaning up %s", err))
			}
		}
	}()

	err = shell.StartContainer()
	if err != nil {
		return fmt.Errorf("cannot start container ", err)
	}

	errs := make(chan error)
	go func() {
		errs <- shell.DockerClient.AttachToContainer(docker.AttachToContainerOptions{
			Container:    container.ID,
			OutputStream: options.OutputStream,
			ErrorStream:  options.ErrorStream,
			InputStream:  options.InputStream,
			Stdout:       true,
			Stdin:        true,
			Stderr:       true,
			Stream:       true,
			RawTerminal:  shell.Tty,
		})
	}()

	if options.Success != nil {
		options.Success <- struct{}{}
	}

	if err != nil {
		return fmt.Errorf("cannot attach to container ", err)
	}

	if options.Detach != nil {
		go func() {
			err := <-options.Detach
			errs <- err
		}()
	}

	myerr := <-errs

	if myerr != nil {
		return fmt.Errorf("attach error %s", err)
	}

	return nil
}

func (shell *Shell) Log(tag string, message string, args ...string) {
	fmt.Printf(fmt.Sprintf("%s %s [%s] %s\n",
		time.Now(),
		shell.LogPrefix,
		tag,
		fmt.Sprintf(message, args)))
	f, err := os.OpenFile("/var/log/manager/manager-shell.log",
		os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	if _, err = f.WriteString(fmt.Sprintf("%s %s [%s] %s\n",
		time.Now(),
		shell.LogPrefix,
		tag,
		fmt.Sprintf(message, args))); err != nil {
		return
	}
}

func (shell *Shell) LogError(err error) {
	shell.Log("error", fmt.Sprintf("%s", err))
}

func (shell *Shell) ResizeTtyTo(id string, height int, width int) {
	shell.DockerClient.ResizeContainerTTY(id, height, width)
}

func (shell *Shell) ResizeTty(id string, isExec bool) {
	height, width := shell.GetTtySize()
	if height == 0 && width == 0 {
		return
	}

	shell.DockerClient.ResizeContainerTTY(id, height, width)
}

func (shell *Shell) MonitorTtySize(id string, isExec bool) error {
	shell.ResizeTty(id, isExec)

	sigchan := make(chan os.Signal, 1)
	gosignal.Notify(sigchan, signal.SIGWINCH)
	go func() {
		for _ = range sigchan {
			shell.ResizeTty(id, isExec)
		}
	}()
	return nil
}

func (shell *Shell) GetTtySize() (int, int) {
	if !shell.Tty {
		return 0, 0
	}
	ws, err := term.GetWinsize(shell.OutFd)
	if err != nil {
		log.Printf("Error getting size: %s", err)
		if ws == nil {
			return 0, 0
		}
	}
	return int(ws.Height), int(ws.Width)
}

func (s *Shell) GetDockerImages() {
	endpoint := "unix:///var/run/docker.sock"
	client, _ := docker.NewClient(endpoint)
	imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})

	s.Images = []string{}
	s.ShellImages = []string{}

	for _, img := range imgs {
		if len(img.RepoTags) > 0 {
			for _, tag := range img.RepoTags {
				if strings.Contains(tag, "manager/") {
					tag = strings.Replace(tag, "manager/", "", 1)
					tag = strings.Replace(tag, ":latest", "", 1)
					s.Images = append(s.Images, tag)

					if strings.Contains(tag, "-shell") {
						s.ShellImages = append(s.ShellImages, tag)
					}
				}
			}
		}
	}
}

func WebSocketShell(c *gin.Context, sharedContext *shared.SharedContext) {
	s := Shell{
		LogPrefix:     "[web shell]",
		SharedContext: sharedContext,
	}

	//environment
	a := models.AccountFromContext(c)
	env := strings.Replace(c.Request.URL.Query().Get("env"), "-base-shell", "", 1)
	fmt.Printf(env)
	out, err := exec.Command("id", "-u", a.Name).Output()

	if err != nil {
		return
	}

	uid := strings.Replace(string(out), "\n", "", 1)

	s.Log("info", "user: %s (%s)", a.Name, uid)

	s.Cmd = []string{"/bin/bash"}
	s.Tty = true

	s.Log("info", "cmd: %s", strings.Join(s.Cmd, " "))

	//go-dockerclient
	endpoint := "unix:///var/run/docker.sock"
	s.DockerClient, err = docker.NewClient(endpoint)

	if err != nil {
		s.LogError(errors.New("cannot connect to docker client"))
		return
	}

	s.GetDockerImages()

	//environment
	s.AccountName = a.Name
	s.AccountUid = uid

	shell_image, err := s.BuildShellImage(env)

	if err != nil {
		s.LogError(err)
		return
	}

	errs := make(chan error)

	r, w := io.Pipe()
	stdinR, stdinW := io.Pipe()

	detach := make(chan error)

	go func() {
		errs <- s.Attach(AttachOptions{
			ShellImage:   shell_image,
			InputStream:  stdinR,
			OutputStream: w,
			ErrorStream:  w,
			Detach:       detach,
		})
	}()

	conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Writer, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	go func(reader io.Reader) {
		for {
			data := make([]byte, 100)
			n, err := reader.Read(data)

			if err != nil {
				errs <- fmt.Errorf("error while reading from docker stream ", err)
				return
			}

			if err := conn.WriteMessage(websocket.TextMessage, data[:n]); err != nil {
				conn.Close()
				errs <- errors.New("error while writing to websocket stream")
				return
			}
		}
	}(r)

	go func(writer io.Writer) {
		for {
			_, p, err := conn.ReadMessage()
			if err != nil {
				errs <- fmt.Errorf("error while reading from websocket ", err)
				return
			}

			writer.Write(p)
		}
	}(stdinW)

	detach <- <-errs
}
