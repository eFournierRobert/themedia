// Package tools is the package containing
// all the request made to the database.
package tools

import (
	"errors"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
)

// GetDb is a function that gets the gorm.DB we need to
// interact with the database. It will get the database auth
// and the database name from the build/.env file.
// It will return a pointer to the gorm.DB or an error if one 
// occured.
func GetDb() (*gorm.DB, error) {
	err := godotenv.Load("build/.env")

	if err != nil {
		log.Fatal("Couldn't read .env file")
		return nil, errors.New("Couldn't read .env file")
	}

	db_username := os.Getenv("MYSQL_USER")
	db_password := os.Getenv("MYSQL_PASSWORD")
	db_name := os.Getenv("MYSQL_DATABASE")

	dsn := db_username + ":" + db_password + "@tcp(127.0.0.1:3306)/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Couldn't connect to database")
		return nil, errors.New("Couldn't connect to database")
	}

	return db, nil
}
