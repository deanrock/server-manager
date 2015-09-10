package controllers

import (
	"../models"
	"../shared"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type SSHPasswordsAPI struct {
	Context *shared.SharedContext
}

func (api *SSHPasswordsAPI) ListPasswords(c *gin.Context) {
	a := models.AccountFromContext(c)
	var passwords []models.SSHPassword
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&passwords)
	c.JSON(200, passwords)
}

func (api *SSHPasswordsAPI) GetPassword(c *gin.Context) {
	a := models.AccountFromContext(c)
	var password models.SSHPassword
	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).Find(&password).Error; err != nil {
		c.String(404, "")
	} else {
		c.JSON(200, password)
	}
}

func (api *SSHPasswordsAPI) validate(c *gin.Context) (*models.SSHPasswordAddForm, *shared.FormErrors) {
	var form models.SSHPasswordAddForm
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Description == "" {
		fe.Add("description", "This field is required.")
	}

	if len(form.Password) <= 10 {
		fe.Add("password", "Password needs to be at least 10 chars long.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *SSHPasswordsAPI) AddPassword(c *gin.Context) {
	a := models.AccountFromContext(c)
	password := models.SSHPassword{}

	password.CreatedAt = time.Now()
	password.UpdatedAt = time.Now()
	password.Added_by_id = c.MustGet("user").(models.User).Id

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	password.Password = form.Password
	password.Description = form.Description
	password.Account_id = a.Id

	api.Context.PersistentDB.Save(&password)
}

func (api *SSHPasswordsAPI) DeletePassword(c *gin.Context) {
	a := models.AccountFromContext(c)
	id := c.Params.ByName("id")

	var password models.SSHPassword

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, id).First(&password).Error; err != nil {
		c.String(404, "")
		return
	}

	api.Context.PersistentDB.Delete(&password)
}
