package controllers

import (
	"regexp"
	"time"

	"../models"
	"../shared"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserForm struct {
	Id          int       `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	First_name  string    `json:"first_name"`
	Last_name   string    `json:"last_name"`
	Email       string    `json:"email"`
	Is_staff    bool      `json:"is_staff"`
	Is_active   bool      `json:"is_active"`
	Last_login  time.Time `json:"last_login"`
	Date_joined time.Time `json:"date_joined"`
}

type UsersAPI struct {
	Context *shared.SharedContext
}

func (api *UsersAPI) ListUsers(c *gin.Context) {
	var users []models.User
	api.Context.PersistentDB.Find(&users)
	c.JSON(200, users)
}

func (api *UsersAPI) GetMyAccess(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	var access models.UserAccess
	api.Context.PersistentDB.Where("user_id = ? AND account_id = ?", user.Id, c.Params.ByName("account")).Find(&access)
	c.JSON(200, access)
}

func (api *UsersAPI) GetUser(c *gin.Context) {
	var user models.User
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	c.JSON(200, user)
}

func (api *UsersAPI) validate(c *gin.Context, exists bool) (*UserForm, *shared.FormErrors) {
	var form UserForm
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Username == "" {
		fe.Add("username", "This field is required.")
	}

	match, _ := regexp.MatchString("^([a-zA-Z][0-9a-zA-Z_]*)$", form.Username)
	if !match {
		fe.Add("username", "Only alphanumeric characters and underscore are allowed.")
	}

	if len(form.Username) > 50 {
		fe.Add("username", "Ensure this value has at most 50 characters.")
	}

	if !exists {
		if len(form.Password) < 8 {
			fe.Add("password", "Ensure this value has at least 8 characters.")
		}
	} else {
		if len(form.Password) > 0 && len(form.Password) < 8 {
			fe.Add("password", "Ensure this value has at least 8 characters.")
		}
	}

	if form.Email == "" {
		fe.Add("email", "This field is required.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *UsersAPI) EditUser(c *gin.Context) {
	var user models.User

	id := c.Params.ByName("id")

	if id != "" {
		if err := api.Context.PersistentDB.Where("id = ?", id).First(&user).Error; err != nil {
			c.String(404, "")
			return
		}
	} else {
		user.Date_joined = time.Now()
	}

	form, fe := api.validate(c, id != "")
	if fe != nil {
		c.JSON(400, fe)
		return
	}

	user.Username = form.Username

	if len(form.Password) > 0 {
		user.Password = models.GeneratePasswordHash(form.Password)
	}

	user.First_name = form.First_name
	user.Last_name = form.Last_name
	user.Email = form.Email
	user.Is_staff = form.Is_staff
	user.Is_active = form.Is_active
	user.Is_superuser = false

	api.Context.PersistentDB.Save(&user)

	c.JSON(200, user)
}

func (api *UsersAPI) GetAccess(c *gin.Context) {
	var access []models.UserAccess
	api.Context.PersistentDB.Where("user_id = ?", c.Params.ByName("id")).Find(&access)
	c.JSON(200, access)
}

func (api *UsersAPI) SetAccess(c *gin.Context) {
	var user models.User
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	var account models.Account
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("account")).First(&account).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	var userAccess models.UserAccess
	if err := api.Context.PersistentDB.Where("user_id = ? AND account_id = ?", user.Id, account.Id).First(&userAccess).Error; err != nil {
		//create
		userAccess = models.UserAccess{
			User_id:     user.Id,
			Account_id:  account.Id,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Added_by_id: models.UserFromContext(c).Id,
		}

		api.Context.PersistentDB.Save(&userAccess)
	}

	var form models.UserAccess
	c.BindWith(&form, binding.JSON)

	userAccess.SshAccess = form.SshAccess
	userAccess.ShellAccess = form.ShellAccess
	userAccess.AppAccess = form.AppAccess
	userAccess.DomainAccess = form.DomainAccess
	userAccess.DatabaseAccess = form.DatabaseAccess
	userAccess.CronjobAccess = form.CronjobAccess
	userAccess.UpdatedAt = time.Now()

	api.Context.PersistentDB.Save(&userAccess)
}

func (api *UsersAPI) RemoveAccess(c *gin.Context) {
	var user models.User
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("id")).First(&user).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}
	var account models.Account
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("account")).First(&account).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	var userAccess models.UserAccess
	if err := api.Context.PersistentDB.Where("user_id = ? AND account_id = ?", user.Id, account.Id).First(&userAccess).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	api.Context.PersistentDB.Delete(&userAccess)
}
