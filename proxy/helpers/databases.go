package helpers

import (
  "../models"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "fmt"
)

func CreateMysqlDatabase(database *models.Database) (bool, error) {

  success := true
  db, err := sql.Open("mysql", "root:"+"password"+"@tcp(127.0.0.1:3306)/")
  defer db.Close()

  if err != nil {
    fmt.Println(err)
    return false, err
  }

  _, err = db.Exec("CREATE DATABASE `"+database.Name+"` CHARACTER SET utf8 COLLATE utf8_general_ci")

  if err != nil {
    return false, err
  }

  _, err = db.Exec("GRANT ALL ON `"+database.Name+"`.* TO `"+database.User+"`@'%%' IDENTIFIED BY '"+database.Password+"'")

  if err != nil {
    return false, err
  }

  _, err = db.Exec("FLUSH PRIVILEGES")

  return success, err
}
