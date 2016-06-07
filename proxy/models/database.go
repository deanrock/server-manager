package models

import (
	"time"
)

type Database struct {
	Id          int       `json:"id"`
	Account_id  int       `json:"account_id"`
	Type        string    `json:"type",gorm:"size:50"`
	Name        string    `json:"name",gorm:"size:50"`
	User        string    `json:"user",gorm:"size:16"`
	Password    string    `json:"password"`
	Added_at    time.Time `json:"added_at"`
	Added_by_id int       `json:"added_by"`
}

func (c Database) TableName() string {
	return "manager_database"
}
