package controllers

import (
	"../models"
	"../shared"
	"errors"
	"github.com/gin-gonic/gin"
)

type TasksAPI struct {
	Context *shared.SharedContext
}

func (api *TasksAPI) ListMyTasks(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var tasks []models.Task

	if user.Is_staff {
		if err := api.Context.PersistentDB.Find(&tasks).Error; err != nil {
			c.AbortWithError(400, errors.New("DB error"))
			return
		}
	} else {
		if err := api.Context.PersistentDB.Where("added_by_id = ?", user.Id).Find(&tasks).Error; err != nil {
			c.AbortWithError(400, errors.New("DB error"))
			return
		}
	}
	c.JSON(200, tasks)
}

func (api *TasksAPI) GetTask(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var task models.Task
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("id")).First(&task).Error; err != nil {
		c.AbortWithError(400, errors.New("DB error"))
		return
	}

	if user.Is_staff || task.Added_by_id == user.Id {
		c.JSON(200, task)
	} else {
		c.AbortWithStatus(401)
	}
}

func (api *TasksAPI) GetTaskLog(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	var task models.Task
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("id")).First(&task).Error; err != nil {
		c.AbortWithError(400, errors.New("DB error"))
		return
	}

	var logs []models.TaskLog
	if err := api.Context.PersistentDB.Where("task_id = ?", c.Params.ByName("id")).Find(&logs).Error; err != nil {
		c.AbortWithError(400, errors.New("DB error"))
		return
	}

	if user.Is_staff || task.Added_by_id == user.Id {
		c.JSON(200, logs)
	} else {
		c.AbortWithStatus(401)
	}
}
