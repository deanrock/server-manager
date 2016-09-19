package shared

import (
	"../realtime"
	"bytes"
	"encoding/json"
	"github.com/fsouza/go-dockerclient"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
)

type Config struct {
	Server_name             string `json:"server_name"`
	Mysql_root_password     string `json:"mysql_root_password"`
	Postgres_root_password  string `json:"postgres_root_password"`
	Mysql_connection_string string `json:"mysql_connection_string"`
}

type SharedContext struct {
	PersistentDB     gorm.DB
	LogDB            gorm.DB
	DockerClient     *docker.Client
	WebsocketHandler realtime.WebsocketHandler
	Config           *Config
}

func (s *SharedContext) InitConfig(path string) {
	c, _ := ioutil.ReadFile(path)
	dec := json.NewDecoder(bytes.NewReader(c))
	var config Config
	dec.Decode(&config)

	s.Config = &config
}

func (s *SharedContext) OpenDB(path string) {
	//sqlite
	db, err := gorm.Open("sqlite3", path)
	if err != nil {
		log.Fatal("database error", err)
	}

	s.PersistentDB = db
}
