package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"../models"
	"../shared"
	//"fmt"
)

type CronJobsAPI struct {
	Context *shared.SharedContext
}

func (api *CronJobsAPI) ListCronjobs(c *gin.Context) {
	a := models.AccountFromContext(c)
	var cronjobs []models.CronJob
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&cronjobs)
	c.JSON(200, cronjobs)
}

func (api *CronJobsAPI) GetCronjob(c *gin.Context) {
	a := models.AccountFromContext(c)
	var cronjob models.CronJob
	if err := api.Context.PersistentDB.Where("account_id = ?, id = ?", a.Id, c.Params.ByName("id")).First(&cronjob).Error; err != nil {
		c.JSON(200, cronjob)
	}else{
		c.String(404, "")
	}
}

func (api *CronJobsAPI) AddCronjob(c *gin.Context) {
	a := models.AccountFromContext(c)

	var form models.CronJobForm
	c.BindWith(&form, binding.MultipartForm)

	fe := shared.FormErrors{}

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	if form.Directory == "" {
		fe.Add("directory", "This field is required.")
	}

	if form.Command == "" {
		fe.Add("command", "This field is required.")
	}

	if form.Cron_expression == "" {
		fe.Add("cron_expression", "This field is required.")
	}

	var image models.Image
	if err := api.Context.PersistentDB.Where("id = ?", form.Image_id).First(&image).Error; err != nil {
		fe.Add("image_id", "Image with this id doesn't exist.")
	}


	if fe.HasErrors() {
		c.JSON(400, fe)
		return
	}

	cx := models.CronJob{
		Name: form.Name,
		Directory: form.Directory,
		Command: form.Command,
		Timeout: form.Timeout,
		Cron_expression: form.Cron_expression,
		Image_id: image.Id,
		Account_id: a.Id,
	}

	api.Context.PersistentDB.Save(&cx)

	c.JSON(200, cx)
}

func (api *CronJobsAPI) EditCronjob(c *gin.Context) {

}
