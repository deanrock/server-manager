package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jinzhu/gorm"

	"./cron"
	"./proxy"
	"./ssh-server"
	"./proxy/models"
	"./proxy/shared"
)

func firstRun() {
	db, err := gorm.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal("database error", err)
	}

	// execute DB schema from Django project
	schema, err := ioutil.ReadFile("schema_v1.sql")
	if err != nil {
		fmt.Print(err)
	}

	db.Exec(string(schema))
}

func createAdminUser() {
	  username := os.Args[2]
		password := os.Args[3]

		if username == "" || password == "" {
			log.Println("Username and password cannot be empty!")
			os.Exit(1)
		}

		// create user
		user := models.User{
			Username: username,
			Password: models.GeneratePasswordHash(password),
			First_name: "Admin",
			Last_name: "Admin",
			Email: "admin@example.com",
			Is_staff: true,
			Is_active: true,
			Is_superuser: true,
		}
		context := shared.SharedContext{}
		context.OpenDB("db.sqlite3")
		context.PersistentDB.Save(&user)
}

func main() {
	log.Println(fmt.Sprintf("starting %s ...", os.Args[1]))

	switch os.Args[1] {
	case "proxy":
		proxy.Start()
	case "cron":
		cron.Start()
	case "ssh":
		ssh_server.Start()
	case "first-run":
		firstRun()
	case "create-admin-user":
		createAdminUser()
	default:
		log.Println("command not found; exiting")
		os.Exit(1)
	}
}
