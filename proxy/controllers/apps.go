package controllers

import (
	"../models"
	"../shared"
	"../tasks"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"regexp"
	"time"
)

type AppsAPI struct {
	Context *shared.SharedContext
}

func (api *AppsAPI) ListApps(c *gin.Context) {
	a := models.AccountFromContext(c)
	var apps []models.App
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&apps)
	c.JSON(200, apps)
}

func (api *AppsAPI) GetApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	api.Context.PersistentDB.Where("app_id = ?", app.Id).Find(&app.Variables)

	c.JSON(200, app)
}

func (api *AppsAPI) validate(c *gin.Context) (*models.App, *shared.FormErrors) {
	var form models.App
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	match, _ := regexp.MatchString("^([a-zA-Z][0-9a-zA-Z_-]*)$", form.Name)
	if !match {
		fe.Add("name", "Only alphanumeric characters, underscore and '-' are allowed.")
	}

	if form.Memory < 0 || form.Memory > 16000 {
		fe.Add("memory", "Valid memory is between 0 and 1600 MB.")
	}

	var images []models.Image
	api.Context.PersistentDB.Find(&images)
	found := false
	for _, image := range images {
		if image.Id == form.Image_id {
			found = true
			break
		}
	}
	if !found {
		fe.Add("image_id", "Image doesn't exist.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *AppsAPI) EditApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	app := models.App{}

	id := c.Params.ByName("id")

	if id != "" {
		if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, id).First(&app).Error; err != nil {
			c.String(404, "")
			return
		}
	} else {
		app.Added_at = time.Now()
		app.Added_by_id = c.MustGet("user").(models.User).Id
		app.Account_id = a.Id
	}

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	app.Name = form.Name
	app.Memory = form.Memory
	app.Image_id = form.Image_id

	api.Context.PersistentDB.Save(&app)

	for _, variable := range form.Variables {
		if variable.Value == "" {
			api.Context.PersistentDB.Where("app_id=? AND name=?", app.Id, variable.Name).Delete(models.AppVariable{})
		} else {
			v := models.AppVariable{}
			if err := api.Context.PersistentDB.Where("app_id=? AND name=?", app.Id, variable.Name).First(&v).Error; err != nil {
				v.Name = variable.Name
				v.App_id = app.Id
			}

			v.Value = variable.Value
			api.Context.PersistentDB.Save(&v)
		}
	}
}

func (api *AppsAPI) StartApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	user := c.MustGet("user").(models.User).Id

	task := tasks.StartApp(&app, a, user, api.Context)

	c.JSON(200, task)
}

func (api *AppsAPI) StopApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	user := c.MustGet("user").(models.User).Id

	task := tasks.StopApp(&app, a, user, api.Context)

	c.JSON(200, task)
}

func (api *AppsAPI) RedeployApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	user := c.MustGet("user").(models.User).Id

	task := tasks.RedeployApp(&app, a, user, api.Context)

	c.JSON(200, task)
}
