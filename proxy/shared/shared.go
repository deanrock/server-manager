package shared

import (
	"../realtime"
	"github.com/fsouza/go-dockerclient"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type SharedContext struct {
	PersistentDB     gorm.DB
	LogDB            gorm.DB
	DockerClient     *docker.Client
	WebsocketHandler realtime.WebsocketHandler
}

func (s *SharedContext) OpenDB(path string) {
	//sqlite
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		log.Fatal("database error", err)
	}

	s.PersistentDB = db
}
