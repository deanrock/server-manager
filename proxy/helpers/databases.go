package helpers

import (
	"../models"
	"../shared"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func CreateMysqlDatabase(sharedContext *shared.SharedContext, database *models.Database) (bool, error) {

	success := true
	db, err := sql.Open("mysql", "root:"+"password"+"@tcp(127.0.0.1:3306)/")
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = db.Exec("CREATE DATABASE `" + database.Name + "` CHARACTER SET utf8 COLLATE utf8_general_ci")

	if err != nil {
		return false, err
	}

	_, err = db.Exec("GRANT ALL ON `" + database.Name + "`.* TO `" + database.User + "`@'%%' IDENTIFIED BY '" + database.Password + "'")

	if err != nil {
		return false, err
	}

	_, err = db.Exec("FLUSH PRIVILEGES")

	return success, err
}

func CreatePostgresDatabase(sharedContext *shared.SharedContext, database *models.Database) (bool, error) {

	success := true
	db, err := sql.Open("postgres", "user=postgres password=password sslmode=disable host=localhost")
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = db.Exec("CREATE DATABASE " + database.Name + " ENCODING 'utf8' LC_COLLATE 'en_US.UTF-8'")

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = db.Exec("CREATE ROLE " + database.User + " WITH LOGIN PASSWORD '" + database.Password + "'")

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = db.Exec("REVOKE ALL PRIVILEGES ON DATABASE " + database.Name + " FROM PUBLIC")

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = db.Exec("GRANT ALL PRIVILEGES ON DATABASE " + database.Name + " TO " + database.User)

	if err != nil {
		fmt.Println(err)
		return false, err
	}

	_, err = db.Exec("FLUSH PRIVILEGES")

	return success, err
}
