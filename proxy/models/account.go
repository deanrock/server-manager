package models

import (
	"../shared"
	"github.com/gin-gonic/gin"
	"time"
)

type Account struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Added_at    time.Time `json:"added_at"`
	Added_by_id int       `json:"added_by"`
	Description string    `json:"description"`

	context *shared.SharedContext
}

func AccountFromContext(c *gin.Context) *Account {
	a := c.MustGet("account").(*Account)
	return a
}

func (a Account) TableName() string {
	return "manager_account"
}

func GetAccountByName(name string, c *shared.SharedContext) *Account {
	var account Account
	if err := c.PersistentDB.Where("name = ?", name).First(&account).Error; err != nil {
		return nil
	}

	account.context = c
	return &account
}

func (a Account) Apps() []App {
	var apps []App
	a.context.PersistentDB.Where("account_id = ?", a.Id).Find(&apps)
	return apps
}
