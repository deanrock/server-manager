package main

import (
	"github.com/robfig/cron"
	"fmt"
	"time"
)

func main() {
	c := cron.New()
	c.AddFunc("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.AddFunc("@hourly",      func() { fmt.Println("Every hour") })
	c.AddFunc("@every 1s", func() { fmt.Println("Every hour thirty") })
	c.Start()

	fmt.Println(c.Entries())

	c.Remove(x)

	for {
		time.Sleep(time.Second*100)
	}
}
