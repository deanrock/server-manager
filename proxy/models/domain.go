package models

import (
	//"../shared"
	"time"
)

type Domain struct {
	Id            int    `json:"id"`
	Account_id    int    `json:"account_id"`
	Name          string `json:"name"`
	Redirect_url  string `json:"redirect_url"`
	Nginx_config  string `json:"nginx_config"`
	Apache_config string `json:"apache_config"`

	Apache_enabled bool `json:"apache_enabled"`
	Ssl_enabled    bool `json:"ssl_enabled"`

	Added_at    time.Time `json:"added_at"`
	Added_by_id int       `json:"added_by"`
}

func (d Domain) TableName() string {
	return "manager_domain"
}
