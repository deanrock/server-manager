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
