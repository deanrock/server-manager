package controllers

import (
	"../models"
	"../shared"
	"../helpers"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"regexp"
	"time"
)

type DatabasesAPI struct {
	Context *shared.SharedContext
}

func (api *DatabasesAPI) ListDatabases(c *gin.Context) {
	a := models.AccountFromContext(c)
	var databases []models.Database
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&databases)
	c.JSON(200, databases)
}

func (api *DatabasesAPI) GetDatabase(c *gin.Context) {
	a := models.AccountFromContext(c)
	var database models.Database

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&database).Error; err != nil {
		c.String(404, "")
		return
	}

	c.JSON(200, database)
}

func (api *DatabasesAPI) validate(c *gin.Context) (*models.Database, *shared.FormErrors) {
	var form models.Database
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Type != "mysql" && form.Type != "postgres" {
		fe.Add("type", "This type is not supported.")
	}

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	match, _ := regexp.MatchString("^([a-zA-Z][0-9a-zA-Z_]*)$", form.Name)
	if !match {
		fe.Add("name", "Only alphanumeric characters and underscore are allowed.")
	}

	if len(form.Name) > 50 {
		fe.Add("name", "Ensure this value has at most 50 characters.")
	}

	if form.User == "" {
		fe.Add("user", "This field is required.")
	}

	match, _ = regexp.MatchString("^([a-zA-Z][0-9a-zA-Z_]*)$", form.User)
	if !match {
		fe.Add("user", "Only alphanumeric characters and underscore are allowed.")
	}

	if len(form.User) > 16 {
		fe.Add("user", "Ensure this value has at most 16 characters.")
	}

	if form.Password == "" {
		fe.Add("password", "This field is required.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *DatabasesAPI) EditDatabase(c *gin.Context) {
	a := models.AccountFromContext(c)
	database := models.Database{}

	id := c.Params.ByName("id")

	if id != "" {
		if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, id).First(&database).Error; err != nil {
			c.String(404, "")
			return
		}
	} else {
		database.Added_at = time.Now()
		database.Added_by_id = c.MustGet("user").(models.User).Id
		database.Account_id = a.Id
	}

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	database.Name = form.Name
	database.User = form.User
	database.Password = form.Password
	database.Type = form.Type

	api.Context.PersistentDB.Save(&database)

	if id == "" {
			success, err := helpers.CreateMysqlDatabase(&database)

			if !success {
				c.JSON(400, err)
				return
			}
	} else {
			// change mysql database password
	}
}
