package controllers

import (
	"github.com/gin-gonic/gin"
	"../models"
	"../shared"
	"errors"
	"time"
	"github.com/gin-gonic/gin/binding"
)

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
        c.Fail(404, errors.New("Not Found"))
        return
    }
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
        c.Fail(404, errors.New("User Not Found"))
        return
    }
	var account models.Account
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("account")).First(&account).Error; err != nil {
        c.Fail(404, errors.New("Account Not Found"))
        return
    }

    var userAccess models.UserAccess
    if err := api.Context.PersistentDB.Where("user_id = ? AND account_id = ?", user.Id, account.Id).First(&userAccess).Error; err != nil {
    	//create
    	userAccess = models.UserAccess{
    		User_id: user.Id,
    		Account_id: account.Id,
    		CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
        c.Fail(404, errors.New("Not Found"))
        return
    }
	var account models.Account
	if err := api.Context.PersistentDB.Where("id = ?", c.Params.ByName("account")).First(&account).Error; err != nil {
        c.Fail(404, errors.New("Not Found"))
        return
    }

    var userAccess models.UserAccess
    if err := api.Context.PersistentDB.Where("user_id = ? AND account_id = ?", user.Id, account.Id).First(&userAccess).Error; err != nil {
    	c.Fail(404, errors.New(" UserAccess doesn't exist"))
    	return
    }

    api.Context.PersistentDB.Delete(&userAccess)
}
