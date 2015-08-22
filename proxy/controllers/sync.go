package controllers

import (
	"../models"
	"../shared"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"log"
	"strings"
	"time"
)

type SyncAPI struct {
	Context *shared.SharedContext
}

func (api *SyncAPI) GetImages(c *gin.Context) {
	images := [...]string{
		"debian7base",
		"debian7basehosting",
		"debian7baseshell",

		"php53-base-hosting",

		"php56-base",
		"php56-base-hosting",
		"php56-base-shell",

		"python27-base",
		"python27-base-shell",

		"python34-base",
		"python34-base-shell",

		"java8-base",
		"java8-base-shell",

		"go1.4-base",
		"go1.4-base-shell",

		"nodejs0.12-base",
		"nodejs0.12-base-shell",

		"ruby22-base",
		"ruby22-base-shell",
	}

	c.JSON(200, images)
}

func (api *SyncAPI) SyncImage(c *gin.Context) {
	name := c.Params.ByName("name")

	//TODO: check that image name only contains A-z0-9-.

	vars, err := json.Marshal(struct {
		ImageName string `json:"image_name"`
	}{
		ImageName: name,
	})

	if err != nil {
		c.AbortWithError(500, err)
	}

	go func(user int) {
		//create task
		task := models.NewTask("sync-image", string(vars), user)
		api.Context.PersistentDB.Save(&task)
		task.NotifyUser(*api.Context, user)

		var success = false
		defer func() {
			task.Duration = time.Now().Sub(task.Added_at).Seconds()
			task.Finished = true
			task.Success = success
			api.Context.PersistentDB.Save(&task)
			task.NotifyUser(*api.Context, user)
		}()

		//build image options
		buf := bytes.NewBuffer(nil)
		opts := docker.BuildImageOptions{
			Name:                fmt.Sprintf("manager/%s", name),
			NoCache:             true,
			RmTmpContainer:      true,
			ForceRmTmpContainer: true,
			OutputStream:        buf,
			RawJSONStream:       true,
			SuppressOutput:      false,
			ContextDir:          fmt.Sprintf("../images/%s/", name),
		}

		//call build image API
		err := api.Context.DockerClient.BuildImage(opts)

		if err != nil {
			log.Println(err)

			l := models.TaskLog{
				TaskId:   task.Id,
				Added_at: time.Now(),
				Value:    fmt.Sprintf("error encountered while building the image: %s", err),
				Type:     "error",
			}

			api.Context.PersistentDB.Save(&l)
			return
		}

		//read output from building the image
		var line = ""
		scanner := bufio.NewScanner(buf)
		for scanner.Scan() {
			line = scanner.Text()

			l := models.TaskLog{
				TaskId:   task.Id,
				Added_at: time.Now(),
				Value:    line,
				Type:     "log",
			}

			api.Context.PersistentDB.Save(&l)

			fmt.Println("line!")
			fmt.Println(line)

			if err != nil {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			log.Println(err)

			l := models.TaskLog{
				TaskId:   task.Id,
				Added_at: time.Now(),
				Value:    fmt.Sprintf("error encountered while reading output: %s", err),
				Type:     "error",
			}

			api.Context.PersistentDB.Save(&l)
			return
		}

		if !strings.Contains(line, "Successfully built") {
			log.Println("Last line doesnt contain Successfully build")

			l := models.TaskLog{
				TaskId:   task.Id,
				Added_at: time.Now(),
				Value:    fmt.Sprintf("last line doesn't contain 'Successfully built'"),
				Type:     "error",
			}

			api.Context.PersistentDB.Save(&l)
			return
		}

		success = true
	}(c.MustGet("user").(models.User).Id)

	c.String(200, "")
}
