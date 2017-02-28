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

func SyncWebServers(user int, sharedContext *shared.SharedContext) models.Task {
	//create task
	task := models.NewTask("sync-web-servers", string("{}"), user)
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

	success = helpers.SyncWebServers(sharedContext, task, nil)

	return task
}

// SyncWebServersAfterRestart re-syncs web servers after restarting the server.
//
// Re-syncing occurs twice:
// - first after starting the server,
// - and then 5 minutes after that (when app containers are hopefully
//   already started).
func SyncWebServersAfterRestart(sharedContext *shared.SharedContext) {
	user := 1 // hopefully this user actually exists :)
	SyncWebServers(user, sharedContext)

	time.Sleep(5 * time.Minute)
	SyncWebServers(user, sharedContext)
}
