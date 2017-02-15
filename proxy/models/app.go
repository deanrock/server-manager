package models

import (
	"fmt"
	"time"

	"../shared"
)

type App struct {
	Id           int       `json:"id"`
	Name         string    `json:"name"`
	Container_id string    `json:"container_id"`
	Memory       int       `json:"memory"`
	Account_id   int       `json:"account_id"`
	Image_id     int       `json:"image_id"`
	Added_at     time.Time `json:"added_at"`
	Added_by_id  int       `json:"added_by"`

	Variables []AppVariable `json:"variables"`
}

func (c App) TableName() string {
	return "manager_app"
}

func (c App) ContainerName(accountName string) string {
	return fmt.Sprintf("app-%s-%s", accountName, c.Name)
}

func (c App) GetVariables(context *shared.SharedContext) []AppVariable {
	var variables []AppVariable
	context.PersistentDB.Where("app_id = ?", c.Id).Find(&variables)
	return variables
}
