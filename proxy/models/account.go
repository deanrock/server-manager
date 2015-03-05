package models

import (
	"time"
)

type Account struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Added_at time.Time `json:"added_at"`
	Added_by_id int `json:"added_by"`
	Description string `json:"description"`
}

func (c Account) TableName() string {
    return "manager_account"
}
