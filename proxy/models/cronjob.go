package models

import (
	"time"
)

type CronJobForm struct {
	Name            string `form:"name"`
	Directory       string `form:"directory"`
	Command         string `form:"command"`
	Timeout         int    `form:"timeout"`
	Cron_expression string `form:"cron_expression"`
	Image           string `form:"image"`
	Enabled         bool   `form:"enabled"`
}

type CronJob struct {
	Id              int       `json:"id"`
	Name            string    `json:"name"`
	Directory       string    `json:"directory"`
	Command         string    `json:"command"`
	Timeout         int       `json:"timeout"`
	Cron_expression string    `json:"cron_expression"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Account_id      int       `json:"account"`
	Added_by_id     int       `json:"added_by"`
	Image           string    `json:"image"`
	Enabled         bool      `json:"enabled"`
	Success         bool      `json:"success"`
}
