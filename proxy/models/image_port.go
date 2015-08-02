package models

import (
	//"time"
)

type ImagePort struct {
	Id int `json:"id"`
	Type string `json:"type"`
	Port string `json:"port"`
	Image_id int `json:"image_id"`
}

func (c ImagePort) TableName() string {
    return "manager_imageport"
}
