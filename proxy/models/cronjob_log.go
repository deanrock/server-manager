package models

import (
	"time"
)

type CronJobLog struct {
	Id        int       `json:"id", sql:"AUTO_INCREMENT"`
	CronJobId int       `json:"cronjob_id"`
	Added_at  time.Time `json:"added_at"`
	Success   bool      `json:"success"`
	Duration  float64   `json:"elapsed_time"`
	Log       string    `json:"log"`
	ExitCode  int       `json:"exit_code"`
}
