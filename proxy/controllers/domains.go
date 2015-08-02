package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"../models"
	"time"
	"../shared"
)

type DomainsAPI struct {
	Context *shared.SharedContext
}

func (api *DomainsAPI) ListDomains(c *gin.Context) {
	a := models.AccountFromContext(c)
	var domains []models.Domain
	api.Context.PersistentDB.Where("account_id = ?", a.Id).Find(&domains)
	c.JSON(200, domains)
}

func (api *DomainsAPI) GetDomain(c *gin.Context) {
	a := models.AccountFromContext(c)
	var domain models.Domain
	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, c.Params.ByName("id")).Find(&domain).Error; err != nil {
		c.String(404, "")
	}else{
		c.JSON(200, domain)
	}
}

func (api *DomainsAPI) validate(c *gin.Context) (*models.Domain, *shared.FormErrors) {
	var form models.Domain
	c.BindWith(&form, binding.JSON)

	fe := shared.NewFormErrors()

	if form.Name == "" {
		fe.Add("name", "This field is required.")
	}

	if fe.HasErrors() {
		return nil, &fe
	}

	return &form, nil
}

func (api *DomainsAPI) EditDomain(c *gin.Context) {
	a := models.AccountFromContext(c)
	domain := models.Domain{}

	id := c.Params.ByName("id")

	if id != "" {
		if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, id).First(&domain).Error; err != nil {
			c.String(404, "")
			return
		}
	}else{
		domain.Added_at = time.Now()
		domain.Added_by_id = c.MustGet("user").(models.User).Id
	}

	form, fe := api.validate(c)

	if fe != nil {
		c.JSON(400, fe)
		return
	}

	domain.Account_id = a.Id
	domain.Name = form.Name
	domain.Redirect_url = form.Redirect_url
	domain.Nginx_config = form.Nginx_config
	domain.Apache_config = form.Apache_config
	domain.Apache_enabled = form.Apache_enabled
	domain.Ssl_enabled = form.Ssl_enabled

	api.Context.PersistentDB.Save(&domain)
}

func (api *DomainsAPI) DeleteDomain(c *gin.Context) {
	a := models.AccountFromContext(c)
	id := c.Params.ByName("id")
	
	var domain models.Domain

	if err := api.Context.PersistentDB.Where("account_id = ? AND id = ?", a.Id, id).First(&domain).Error; err != nil {
		c.String(404, "")
		return
	}

	api.Context.PersistentDB.Delete(&domain)
}
