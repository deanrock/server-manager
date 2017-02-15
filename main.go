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
	default:
		log.Println("command not found; exiting")
		os.Exit(1)
	}
}
