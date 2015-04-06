package controllers

import (
	"github.com/robfig/cron"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"../models"
	"time"
	"../shared"
	"fmt"
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
	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&cronjob).Error; err != nil {
		c.String(404, "")
	}else{
		c.JSON(200, cronjob)
	}
}

func (api *CronJobsAPI) validate(c *gin.Context) (*models.CronJobForm, *shared.FormErrors) {
	var form models.CronJobForm
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	if form.Directory == "" {
		fe.Add("directory", "This field is required.")
	}

	if form.Command == "" {
		fe.Add("command", "This field is required.")
	}

	if form.Timeout <= 0 {
		fe.Add("timeout", "Timeout must be more than 0 seconds.")
	}else if form.Timeout > 3600 {
		fe.Add("timeout", "Timeout must be less than or equal 3600 seconds.")
	}

	if form.Cron_expression == "" {
		fe.Add("cron_expression", "This field is required.")
	}else{
		if _, err := cron.Parse(form.Cron_expression); err != nil {
			fe.Add("cron_expression", fmt.Sprintf("Invalid cron expression: %s.", err))
		}
	}

	if form.Image == "" {
		fe.Add("image", "This field is required.")
	}

	//TODO: check if image exists and it's valid

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *CronJobsAPI) AddCronjob(c *gin.Context) {
	a := models.AccountFromContext(c)

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	cx := models.CronJob{
		Name: form.Name,
		Directory: form.Directory,
		Command: form.Command,
		Timeout: form.Timeout,
		Cron_expression: form.Cron_expression,
		Image: form.Image,
		Account_id: a.Id,
		Enabled: form.Enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	api.Context.PersistentDB.Save(&cx)

	c.JSON(200, cx)
}

func (api *CronJobsAPI) EditCronjob(c *gin.Context) {
	a := models.AccountFromContext(c)

	var cronjob models.CronJob
	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&cronjob).Error; err != nil {
		c.String(404, "")
		return
	}

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	cronjob.Name = form.Name
	cronjob.Directory = form.Directory
	cronjob.Command = form.Command
	cronjob.Timeout = form.Timeout
	cronjob.Cron_expression = form.Cron_expression
	cronjob.Image = form.Image
	cronjob.Account_id = a.Id
	cronjob.Enabled = form.Enabled
	cronjob.UpdatedAt = time.Now()

	api.Context.PersistentDB.Save(&cronjob)

	c.JSON(200, cronjob)
}
