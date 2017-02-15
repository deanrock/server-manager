package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"../container"
	"../helpers"
	"../models"
	"../shared"
	"github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
)

type SyncAPI struct {
	Context *shared.SharedContext
}

func (api *SyncAPI) GetImages(c *gin.Context) {
	images := [...]string{
		"debian7base",
		"debian7basehosting",

		"debian8base",

		"php53-base-hosting",

		"php56-base",
		"php56-base-hosting",
		"php56-base-shell",

		"php7-base",

		"python27-base",
		"python27-base-shell",

		"python34-base-shell",
		"python35-base-shell",

		"java8-base",
		"java8-base-shell",

		"go1.4-base",
		"go1.4-base-shell",

		"ruby22-base",
		"ruby22-base-shell",

		"nodejs0.12-base",
		"nodejs0.12-base-shell",

		"nodejs6-base",
		"nodejs6-base-shell",

		"nodejs4-base",
		"nodejs4-base-shell",

		"mongo3.2-base",

		"elixir1.3-base",
		"elixir1.3-base-shell",

		"hhvm-base",

		"memcached-base",

		"rabbitmq-base",

		"redis-base",
	}

	c.JSON(200, images)
}

func (api *SyncAPI) SyncImage(c *gin.Context) {
	name := c.Params.ByName("name")

	no_cache := false
	if c.Query("no-cache") == "true" {
		no_cache = true
	}

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
			NoCache:             no_cache,
			RmTmpContainer:      true,
			ForceRmTmpContainer: true,
			OutputStream:        buf,
			RawJSONStream:       true,
			SuppressOutput:      false,
			ContextDir:          fmt.Sprintf("../images/%s/", name),
		}

		task.Log(fmt.Sprintf("using no-cache: %t", no_cache), "info", api.Context)

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

		err = container.ReadOutputFromBuildImage(api.Context, task, buf)
		if err != nil {
			return
		}

		success = true
	}(c.MustGet("user").(models.User).Id)

	c.String(200, "")
}

func (api *SyncAPI) SyncWebServers(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	//create task
	task := models.NewTask("sync-web-servers", string("{}"), user)
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

	success = helpers.SyncWebServers(api.Context, task, nil)

	if success {
		c.String(200, "")
	}
}
