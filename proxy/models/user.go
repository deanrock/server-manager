package models

import (
	"../shared"
)

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
}

func (u User) TableName() string {
	return "auth_user"
}

func FindUserById(c *shared.SharedContext, id int) (*User, error) {
	var user User
	if err := c.PersistentDB.Where("id = ?", id).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
