package models

import (
	"../shared"
	"time"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id int `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"`
	First_name string `json:"first_name"`
	Last_name string `json:"last_name"`
	Email string `json:"email"`
	Is_staff bool `json:"is_staff"`
	Last_login time.Time `json:"last_login"`
	Date_joined time.Time `json:"date_joined"`
}

func UserFromContext(c *gin.Context) User {
	a := c.MustGet("user").(User)
	return a
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
