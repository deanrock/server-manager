package models

import (
	"time"
)

type TaskLog struct {
	Id       int       `json:"id"`
	TaskId   int       `json:"task_id"`
	Added_at time.Time `json:"added_at"`
	Value    string    `json:"value"`
	Type     string    `json:"type"`
}
