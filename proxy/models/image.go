package models

import (
	"../shared"
	"time"
)

type Image struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Added_at    time.Time `json:"added_at"`

	Ports     []ImagePort     `json:"ports"`
	Variables []ImageVariable `json:"variables"`
}

func (c Image) TableName() string {
	return "manager_image"
}

func (c Image) GetVariables(context *shared.SharedContext) []ImageVariable {
	var variables []ImageVariable
	context.PersistentDB.Where("image_id = ?", c.Id).Find(&variables)
	return variables
}
