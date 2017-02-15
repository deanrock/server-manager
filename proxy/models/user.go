package models

import (
	"crypto/aes"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"../shared"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/pbkdf2"
)

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	Password     string    `json:"-"`
	First_name   string    `json:"first_name"`
	Last_name    string    `json:"last_name"`
	Email        string    `json:"email"`
	Is_staff     bool      `json:"is_staff"`
	Is_superuser bool      `json:"is_superuser"`
	Is_active    bool      `json:"is_active"`
	Last_login   time.Time `json:"last_login"`
	Date_joined  time.Time `json:"date_joined"`
}

func UserFromContext(c *gin.Context) User {
	a := c.MustGet("user").(User)
	return a
}

func (u User) TableName() string {
	return "auth_user"
}

func FindUserById(c *shared.SharedContext, id int) (*User, error) {
	var user User
	if err := c.PersistentDB.Where("id = ?", id).Find(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func CheckPassword(passwordHash, rawPassword string) error {
	parts := strings.Split(passwordHash, "$")
	if len(parts) != 4 {
		return errors.New("Wrong password hash format")
	}

	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("Cannot convert iterations to integer")
	}
	salt := parts[2]
	hash := []byte(parts[3])

	key := []byte(base64.StdEncoding.EncodeToString(pbkdf2.Key([]byte(rawPassword), []byte(salt), iterations, 32, sha256.New)))
	if len(key) != len(hash) {
		return errors.New("Key and hash lengths don't match")
	}
	if subtle.ConstantTimeCompare(key, hash) != 1 {
		return errors.New("Password doesnt match")
	}

	return nil
}

func GeneratePasswordHash(rawPassword string) string {
	saltlen := 12

	// generate salt
	header := make([]byte, saltlen+aes.BlockSize)
	salt := header[:saltlen]
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		panic(err)
	}

	saltString := base64.StdEncoding.EncodeToString(salt)
	iterations := 12000
	key := base64.StdEncoding.EncodeToString(pbkdf2.Key([]byte(rawPassword), []byte(saltString), iterations, 32, sha256.New))

	return fmt.Sprintf("pbkdf2_sha256$%d$%s$%s", iterations, saltString, key)
}
