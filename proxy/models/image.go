package models

import (
	"time"
)

type Image struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Type string `json:"type"`
	Added_at time.Time `json:"added_at"`

	Ports []ImagePort `json:"ports"`
	Variables []ImageVariable `json:"variables"`
}

func (c Image) TableName() string {
    return "manager_image"
}
