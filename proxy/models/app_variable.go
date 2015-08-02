package models

import (
)

type AppVariable struct {
	Id int `json:"id"`
	App_id int `json:"app_id"`
	Name string `json:"name"`
	Value string `json:"value"`
}

func (c AppVariable) TableName() string {
    return "manager_appimagevariable"
}
