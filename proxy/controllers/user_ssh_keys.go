package controllers

import (
	"../models"
	"../shared"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"time"
)

type UserSSHKeysAPI struct {
	Context *shared.SharedContext
}

func (api *UserSSHKeysAPI) ListKeys(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	var keys []models.UserSSHKey
	api.Context.PersistentDB.Where("user_id = ?", user).Find(&keys)
	c.JSON(200, keys)
}

func (api *UserSSHKeysAPI) GetKey(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	var key models.UserSSHKey
	if err := api.Context.PersistentDB.Where("user_id = ? AND id = ?", user, c.Params.ByName("id")).Find(&key).Error; err != nil {
		c.String(404, "")
	} else {
		c.JSON(200, key)
	}
}

func (api *UserSSHKeysAPI) validate(c *gin.Context) (*models.UserSSHKey, *shared.FormErrors) {
	var form models.UserSSHKey
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	if form.SSHKey == "" {
		fe.Add("ssh_key", "This field is required.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *UserSSHKeysAPI) EditKey(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	var key models.UserSSHKey

	id := c.Params.ByName("id")

	if id != "" {
		if err := api.Context.PersistentDB.Where("user_id = ? AND id = ?", user, id).First(&key).Error; err != nil {
			c.String(404, "")
			return
		}
	} else {
		key.Added_at = time.Now()
		key.Added_by_id = user
	}

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	key.Name = form.Name
	key.SSHKey = form.SSHKey
	key.User_id = user

	api.Context.PersistentDB.Save(&key)
}

func (api *UserSSHKeysAPI) DeleteKey(c *gin.Context) {
	user := c.MustGet("user").(models.User).Id

	var key models.UserSSHKey

	id := c.Params.ByName("id")

	if err := api.Context.PersistentDB.Where("user_id = ? AND id = ?", user, id).First(&key).Error; err != nil {
		c.String(404, "")
		return
	}

	api.Context.PersistentDB.Delete(&key)
}
