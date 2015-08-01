package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"../shared"
)

type AccountsAPI struct {
	Context *shared.SharedContext
}

func (api *AccountsAPI) ListAccounts(c *gin.Context) {
	var accounts []models.Account
	api.Context.PersistentDB.Find(&accounts)

	var allowed []models.Account
	userAccess := c.MustGet("userAccess").([]models.UserAccess)

	for _, a:= range(accounts) {
		for _, ua := range(userAccess) {
			if a.Id == ua.Account_id {
				allowed = append(allowed, a)
			}
		}
	}
	c.JSON(200, allowed)
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
