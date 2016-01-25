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

func GetImages(context *shared.SharedContext) []Image {
	var images []Image
	context.PersistentDB.Find(&images)

	var ports []ImagePort
	context.PersistentDB.Find(&ports)

	for k, i := range images {
		for _, v := range ports {
			if v.Image_id == i.Id {
				images[k].Ports = append(images[k].Ports, v)
			}
		}
	}

	var variables []ImageVariable
	context.PersistentDB.Find(&variables)

	for k, i := range images {
		for _, v := range variables {
			if v.Image_id == i.Id {
				images[k].Variables = append(images[k].Variables, v)
			}
		}
	}

	return images
}
