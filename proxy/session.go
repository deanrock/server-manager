package proxy

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func openSqliteDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return db
}

type djangoSession struct {
	UserId int `json:"_auth_user_id"`
}

func readCookie(c *gin.Context) *int {
	sessionID, err := c.Request.Cookie("sessionid")

	if err != nil {
		log.Println(err)
		return nil
	}

	db := openSqliteDatabase()
	defer db.Close()

	var data string
	err = db.QueryRow("SELECT session_data FROM django_session WHERE session_key = ?", sessionID.Value).Scan(&data)

	if err != nil {
		log.Println("Database error", err)
		return nil
	}

	text, err := base64.StdEncoding.DecodeString(data)

	if err != nil {
		log.Println("base64 ", err)
		return nil
	}

	s := string(text)

	i := strings.Index(s, ":")

	if i < 0 {
		log.Println("i < 0")
		return nil
	}

	s = s[i+1:]

	var session djangoSession
	err = json.Unmarshal([]byte(s), &session)

	if err != nil {
		log.Println(err)
		return nil
	}

	if session.UserId <= 0 {
		return nil
	}

	return &session.UserId
}
