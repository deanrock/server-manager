package controllers

import (
	"fmt"
	"regexp"
	"time"

	"../helpers"
	"../models"
	"../shared"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AccountsAPI struct {
	Context *shared.SharedContext
}

func (api *AccountsAPI) ListAccounts(c *gin.Context) {
	var accounts []models.Account
	api.Context.PersistentDB.Find(&accounts)

	var allowed []models.Account
	userAccess := c.MustGet("userAccess").([]models.UserAccess)

	for _, a := range accounts {
		for _, ua := range userAccess {
			if a.Id == ua.Account_id {
				allowed = append(allowed, a)
			}
		}
	}
	c.JSON(200, allowed)
}

func (api *AccountsAPI) validate(c *gin.Context) (*models.Account, *shared.FormErrors) {
	var form models.Account
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	match, _ := regexp.MatchString("^([a-z][0-9a-z_-]*)$", form.Name)
	if !match {
		fe.Add("name", "Only lowercase alphanumeric characters, underscore and '-' are allowed.")
	}

	if !fe.HasErrors() {
		err := helpers.SyncAccount(form.Name, false)
		if err != nil {
			fe.Add("name", fmt.Sprintf("Cannot sync the account: %s", err))
		}
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *AccountsAPI) AddAccount(c *gin.Context) {
	account := models.Account{}

	account.Added_at = time.Now()
	account.Added_by_id = c.MustGet("user").(models.User).Id

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	account.Name = form.Name
	account.Description = form.Description

	api.Context.PersistentDB.Save(&account)
}

func (api *AccountsAPI) ListAllAccounts(c *gin.Context) {
	var accounts []models.Account
	api.Context.PersistentDB.Find(&accounts)
	c.JSON(200, accounts)
}

func (api *AccountsAPI) GetAccountByName(c *gin.Context) {
	account := models.GetAccountByName(c.Params.ByName("name"), api.Context)
	c.JSON(200, account)
}

func (api *AccountsAPI) GetApps(c *gin.Context) {
	account := models.GetAccountByName(c.Params.ByName("name"), api.Context)
	c.JSON(200, account.Apps())
}
