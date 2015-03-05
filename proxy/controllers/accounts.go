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
	var account models.Account
	api.Context.PersistentDB.Where("name = ?", c.Params.ByName("name")).First(&account)
	c.JSON(200, account)
}
