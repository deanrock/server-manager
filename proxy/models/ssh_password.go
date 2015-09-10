package models

import (
	"time"
)

type SSHPassword struct {
	Id          int       `json:"id"`
	Account_id  int       `json:"account_id"`
	Password    string    `json:"-"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Added_by_id int       `json:"added_by"`
}

type SSHPasswordAddForm struct {
	Id          int       `json:"id"`
	Account_id  int       `json:"account_id"`
	Password    string    `json:"password"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Added_by_id int       `json:"added_by"`
}
