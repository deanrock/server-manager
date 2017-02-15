package controllers

import (
	"time"

	"../helpers"
	"../models"
	"../shared"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserLoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthAPI struct {
	Context *shared.SharedContext
}

func (api *AuthAPI) Login(c *gin.Context) {
	var form UserLoginForm
	c.BindWith(&form, binding.JSON)

	var user models.User
	if err := api.Context.PersistentDB.Where("username = ?", form.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(400)
		return
	}

	if models.CheckPassword(user.Password, form.Password) != nil {
		c.AbortWithStatus(400)
		return
	}

	if !user.Is_active {
		c.AbortWithStatus(400)
	}

	user.Last_login = time.Now()
	api.Context.PersistentDB.Save(&user)

	helpers.SetSessionValue(c, "user_id", user.Id)
}
