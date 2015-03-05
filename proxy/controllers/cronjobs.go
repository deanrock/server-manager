package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"../shared"
)

type CronJobsAPI struct {
	Context *shared.SharedContext
}

func (api *CronJobsAPI) ListCronjobs(c *gin.Context) {
	var cronjobs []models.CronJob
	api.Context.PersistentDB.Find(&cronjobs)
	c.JSON(200, cronjobs)
}

func (api *CronJobsAPI) GetCronjob(c *gin.Context) {
	
}

func (api *CronJobsAPI) AddCronjob(c *gin.Context) {
	cx := models.CronJob{
		Name: "yo",
	}

	api.Context.PersistentDB.Save(&cx)
}

func (api *CronJobsAPI) EditCronjob(c *gin.Context) {

}
