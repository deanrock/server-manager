package models

import (
	"time"
	"../shared"
)

type Account struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Added_at time.Time `json:"added_at"`
	Added_by_id int `json:"added_by"`
	Description string `json:"description"`

	context *shared.SharedContext
}

func (a Account) TableName() string {
    return "manager_account"
}

func GetAccountByName(name string, c *shared.SharedContext) *Account {
	var account Account
	c.PersistentDB.Where("name = ?", name).First(&account)

	account.context = c
	return &account
}

func (a Account) Apps() []App {
	var apps []App
	a.context.PersistentDB.Where("account_id = ?", a.Id).Find(&apps)
	return apps
}
