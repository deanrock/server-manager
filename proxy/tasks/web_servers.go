package tasks

import (
	"../helpers"
	"../models"
	"../shared"
	"encoding/json"
	"time"
)

func SyncWebServersForAccount(account *models.Account, user int, sharedContext *shared.SharedContext) models.Task {
	data, _ := json.Marshal(struct {
		Account models.Account `json:"account"`
	}{
		Account: *account,
	})

	//create task
	task := models.NewTask("sync-web-servers-for-domain", string(data), user)
	sharedContext.PersistentDB.Save(&task)
	task.NotifyUser(*sharedContext, user)

	var success = false
	defer func() {
		task.Duration = time.Now().Sub(task.Added_at).Seconds()
		task.Finished = true
		task.Success = success
		sharedContext.PersistentDB.Save(&task)
		task.NotifyUser(*sharedContext, user)
	}()

	success = helpers.SyncWebServers(sharedContext, task, account)

	return task
}
