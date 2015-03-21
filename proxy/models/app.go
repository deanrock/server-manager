package models

import (
	"time"
)

type App struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Container_id string `json:"container_id"`
	Memory int `json:"memory"`
	Account_id int `json:"account_id"`
	Image_id int `json:"image_id"`
	Added_at time.Time `json:"added_at"`
	Added_by_id int `json:"added_by"`
}

func (c App) TableName() string {
    return "manager_app"
}
