package models

import (
	//"time"
)

type ImageVariable struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Default string `json:"default"`
	Filename string `json:"filename"`
	Image_id int `json:"image_id"`
}

func (c ImageVariable) TableName() string {
    return "manager_imagevariable"
}
