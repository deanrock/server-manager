package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/wader/gormstore"
)

func GetSessionValue(c *gin.Context, name string) (interface{}, error) {
	s, _ := c.Get("store")
	store := s.(*gormstore.Store)
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return "", errors.New("session error")
	} else {
		return session.Values[name], nil
	}
}

func SetSessionValue(c *gin.Context, name string, value interface{}) error {
	s, _ := c.Get("store")
	store := s.(*gormstore.Store)
	session, err := store.Get(c.Request, "session")
	if err != nil {
		return errors.New("session error")
	} else {
		session.Values[name] = value
		store.Save(c.Request, c.Writer, session)
	}

	return nil
}
