package controllers

import (
	"github.com/gin-gonic/gin"
	//"github.com/gin-gonic/gin/binding"
	"../models"
	//"time"
	"../shared"
)

type AppsAPI struct {
	Context *shared.SharedContext
}

func (api *AppsAPI) ListApps(c *gin.Context) {
	a := models.AccountFromContext(c)
	var apps []models.App
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&apps)
	c.JSON(200, apps)
}

func (api *AppsAPI) GetApp(c *gin.Context) {
	a := models.AccountFromContext(c)
	var app models.App
	
	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}

	c.JSON(200, app)
}

func (api *AppsAPI) GetAppVariables(c *gin.Context) {
	a := models.AccountFromContext(c)

	var app models.App
	var vars []models.AppVariable

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).First(&app).Error; err != nil {
		c.String(404, "")
		return
	}
	
	api.Context.PersistentDB.Where("app_id = ?", app.Id).Find(&vars)

	c.JSON(200, vars)
}
