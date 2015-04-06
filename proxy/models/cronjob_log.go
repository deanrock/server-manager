package models

import (
	"time"
)

type CronJobLog struct {
	ID           int         `json:"id", sql:"AUTO_INCREMENT"`
	CronJobID    int         `json:"cronjob_id"`
	CreatedAt    time.Time   `json:"created_at"`
	Success      bool        `json:"success"`
	ElapsedTime  int         `json:"elapsed_time"`
	Log          string      `json:"log"`
	ExitCode     int         `json:"exit_code"`
	Metadata     string      `json:"metadata"`
}
