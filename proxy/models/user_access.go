package models

import (
	//"../shared"
	"time"
)

type UserAccess struct {
	Id int `json:"id"`
	User_id int `json:"user_id"`
	Account_id int `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Added_by_id int `json:"added_by"`

	SshAccess bool `json:"ssh_access"`
	ShellAccess bool `json:"shell_access"`
	AppAccess bool `json:"app_access"`
	DomainAccess bool `json:"domain_access"`
	DatabaseAccess bool `json:"database_access"`
	CronjobAccess bool `json:"cronjob_access"`
}
