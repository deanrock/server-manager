package models

import (
	"time"
)

type CronJob struct {
	Id int
	Name string
	Directory string
	Script_file string
	Timeout int
	Cron_expression string
	Added_at time.Time
	Account_id int
	Added_by_id int
	Image_id int
}

func (c CronJob) TableName() string {
    return "manager_cronjob"
}

//("id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
//"name" varchar(255) NOT NULL,
//"directory" varchar(255) NOT NULL, 
//"script_file" varchar(255) NOT NULL,
//"timeout" integer NOT NULL,
//"cron_expression" varchar(255) NOT NULL,
//"added_at" datetime NULL,
//"account_id" integer NOT NULL REFERENCES "manager_account" ("id"),
//"added_by_id" integer NOT NULL REFERENCES "auth_user" ("id"),
//"image_id" integer NOT NULL REFERENCES "manager_image" ("id"));
