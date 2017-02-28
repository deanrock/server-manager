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
	"../tasks"
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

func (api *SyncAPI) PullImage(c *gin.Context) {
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

		//pull image options
		buf := bytes.NewBuffer(nil)
		opts := docker.PullImageOptions{
			Repository:    "deanrock/server-manager",
			Tag:           name,
			Registry:      "",
			OutputStream:  buf,
			RawJSONStream: true,
		}

		task.Log(fmt.Sprintf("pulling image"), "info", api.Context)

		//call pull image API
		err := api.Context.DockerClient.PullImage(opts, docker.AuthConfiguration{})

		if err != nil {
			log.Println(err)

			l := models.TaskLog{
				TaskId:   task.Id,
				Added_at: time.Now(),
				Value:    fmt.Sprintf("error encountered while pulling the image: %s", err),
				Type:     "error",
			}

			api.Context.PersistentDB.Save(&l)
			return
		}

		err = container.ReadOutputFromPullImage(api.Context, task, buf)
		if err != nil {
			return
		}

		success = true
	}(c.MustGet("user").(models.User).Id)

	c.String(200, "")
}

func (api *SyncAPI) SyncWebServers(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	task := tasks.SyncWebServers(user, api.Context)

	if task.Success {
		c.String(200, "")
	}
}

func (api *SyncAPI) SyncAccounts(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	//create task
	task := models.NewTask("sync-users", string("{}"), user)
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

	var accounts []models.Account
	if err := api.Context.PersistentDB.Find(&accounts).Error; err != nil {
		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    fmt.Sprintf("error encountered while fetching accounts: %s", err),
			Type:     "error",
		}

		api.Context.PersistentDB.Save(&l)
		return
	}

	for _, account := range accounts {
		l := models.TaskLog{
			TaskId:   task.Id,
			Added_at: time.Now(),
			Value:    fmt.Sprintf("syncing account: %s", account.Name),
			Type:     "log",
		}

		api.Context.PersistentDB.Save(&l)

		err := helpers.SyncAccount(account.Name, true)
		if err != nil {
			l := models.TaskLog{
				TaskId:   task.Id,
				Added_at: time.Now(),
				Value:    fmt.Sprintf("cannot sync account %s: %s", account.Name, err),
				Type:     "error",
			}

			api.Context.PersistentDB.Save(&l)
			return
		}
	}

	success = true
}

func (api *SyncAPI) PurgeOldLogs(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	//create task
	task := models.NewTask("purge-old-logs", string("{}"), user)
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

	// delete logs older than 7 days
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	// delete cronjob logs
	api.Context.PersistentDB.Where("added_at < ?", sevenDaysAgo).Delete(&models.CronJobLog{})

	// delete tasks
	api.Context.PersistentDB.Where("added_at < ?", sevenDaysAgo).Delete(&models.Task{})

	// get first available task
	var firstTask models.Task
	if err := api.Context.PersistentDB.Order("id").First(&firstTask).Error; err == nil {
		// we got a task, delete task log for every task before this one
		api.Context.PersistentDB.Where("task_id < ?", firstTask.Id).Delete(&models.TaskLog{})
	}

	success = true
	c.JSON(200, gin.H{})
}
