package models

import (
	"../shared"
	"encoding/json"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Variables   string    `json:"variables"`
	Added_at    time.Time `json:"added_at"`
	Added_by_id int       `json:"added_by"`
	Finished    bool      `json:"finished"`
	Duration    float64   `json:"duration"`
	Success     bool      `json:"success"`
	Status      string    `json:"status"`
	Percent     int       `json:"percent"`
	Account_id  *int      `json:"account_id"` //if it belongs to an account
}

func (t Task) Log(message string, messageType string, context *shared.SharedContext) {
	l := TaskLog{
		TaskId:   t.Id,
		Added_at: time.Now(),
		Value:    message,
		Type:     messageType,
	}

	context.PersistentDB.Save(&l)
}

func (t Task) NotifyUser(c shared.SharedContext, uid int) {
	j, _ := json.Marshal(struct {
		Type string `json:"type"`
		Task Task   `json:"task"`
	}{
		Type: "update-task",
		Task: t,
	})
	c.WebsocketHandler.SendToUser(j, uid)
}

func NewTask(name string, variables string, added_by_id int) Task {
	task := Task{
		Name:        name,
		Added_at:    time.Now(),
		Added_by_id: added_by_id,
		Variables:   variables,
	}

	return task
}

func RunningTasksForUser(c *shared.SharedContext, user User) []Task {
	var tasks []Task
	if err := c.PersistentDB.Where("added_by_id = ? AND finished != 1", user.Id).Find(&tasks).Error; err != nil {
		return tasks
	}

	return tasks
}

func CancelAllTasks(c *shared.SharedContext) {
	var tasks []Task
	if err := c.PersistentDB.Where("finished != 1").Find(&tasks).Error; err != nil {

	}

	for _, task := range tasks {
		task.Success = false
		task.Finished = true
		task.Status = "killed"
		c.PersistentDB.Save(&task)
	}
}
