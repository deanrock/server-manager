package main

import (
	//"bufio"
	"net/http"
	//	"fmt"
	"github.com/fsouza/go-dockerclient"
	"io"
	"os/exec"
	"strings"
	"log"
	//"os"
	//"time"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Shell struct {
	Images []string
	ShellImages []string
}

func (s *Shell) getDockerImages() {
	endpoint := "unix:///var/run/docker.sock"
    client, _ := docker.NewClient(endpoint)
    imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})

    s.Images = []string{}
    s.ShellImages = []string{}

    for _, img := range imgs {
        if len(img.RepoTags) > 0 {
        	tag := img.RepoTags[0]
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

func (s *Shell) containerAttachHandler(c *gin.Context) {
	//go-dockerclient
	endpoint := "unix:///var/run/docker.sock"
	dockerClient, _ := docker.NewClient(endpoint)

	//environment
	account := c.Params.ByName("account")
	env := c.Request.URL.Query().Get("env")

	found := false
	for _, e := range s.ShellImages {
		if e == env {
			found = true
		}
	}
	
	if !found {
		return
	}

	out, err := exec.Command("id","-u",account).Output()

	container, err := dockerClient.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			OpenStdin: true,
			Tty: true,
			Cmd:   []string{"/mystuff/start.sh"},
			Image: "manager/"+env,
			Env: []string{"USER="+account, "USERID="+strings.Replace(string(out), "\n", "", 1)},
		},
	})

	defer func() {
		//cleanup
    	log.Printf("cleanup shell %s %s", account, env)
    	err := dockerClient.RemoveContainer(docker.RemoveContainerOptions{
    		ID: container.ID,
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

	err = dockerClient.StartContainer(container.ID,
		&docker.HostConfig{
			Binds: []string{"/home/"+account+":/home/"+account},
			ExtraHosts: []string{"mysql:172.17.42.1"},
		})

	if err != nil {
		log.Printf("cannot start container ", err)
		return
	}

    conn, err := websocket.Upgrade(c.Writer, c.Request, nil, 1024, 1024)
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(c.Writer, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        log.Println(err)
        return
    }

	r, w := io.Pipe()
	stdinR, stdinW := io.Pipe()

	go dockerClient.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: w,
		ErrorStream:  w,
		InputStream:  stdinR,
		Stdout:       true,
		Stdin:        true,
		Stderr:       true,
		Stream:       true,
		RawTerminal: true,
	})
	if err != nil {
		log.Printf("cannot attach to container ", err)
		return
	}

	go func(reader io.Reader) {
		for {
			data := make([]byte, 100)
			n, err := reader.Read(data)

			if err != nil {
				log.Printf("error while reading from docker stream ", err)
			}

			 if err := conn.WriteMessage(websocket.TextMessage, data[:n]); err != nil {
                conn.Close()
                break
            }
		}
	}(r)

	for {
        _, p, err := conn.ReadMessage()
        if err != nil {
            log.Printf("error while reading from websocket ", err)
            return
        }

        stdinW.Write(p)
    }
}
