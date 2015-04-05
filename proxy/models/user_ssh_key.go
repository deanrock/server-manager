package models

import (
	"time"
	"../shared"
)

type UserSSHKey struct {
	Id int `json:"id"`
	User_id int `json:"user_id"`
	Name string `json:"name"`
	SSHKey string `json:"ssh_key"`
	Added_at time.Time `json:"added_at"`
	Added_by_id int `json:"added_by"`
	Description string `json:"description"`
}

func (u UserSSHKey) TableName() string {
    return "manager_usersshkey"
}

func GetAllUserSSHKeys(c *shared.SharedContext) []UserSSHKey {
	var userSSHKey []UserSSHKey
	c.PersistentDB.Find(&userSSHKey)
	return userSSHKey
}

func FindUserSSHKeyByKey(c *shared.SharedContext, key string) (*UserSSHKey, error) {
	var userSSHKey UserSSHKey
	if err := c.PersistentDB.Where("ssh_key = ?", key).Find(&userSSHKey).Error; err != nil {
		return nil, err
	}

	return &userSSHKey, nil
}
