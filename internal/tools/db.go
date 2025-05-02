// Package tools is the package containing
// all the request made to the database.
package tools

import (
	"errors"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// GetDb is a function that gets the gorm.DB we need to
// interact with the database. It will get the database auth
// and the database name from the build/.env file.
// It will return an error if one occured.
func GetDb() error {
	db_username := os.Getenv("MYSQL_USER")
	db_password := os.Getenv("MYSQL_PASSWORD")
	db_name := os.Getenv("MYSQL_DATABASE")

	dsn := db_username + ":" + db_password + "@tcp(127.0.0.1:3306)/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Couldn't connect to database")
		return errors.New("Couldn't connect to database")
	}

	DB = db
	return nil
}
