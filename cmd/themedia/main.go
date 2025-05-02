// Package main is the main package of the API.
// The app starts here.
package main

import (
	"fmt"

	"github.com/eFournierRobert/themedia/internal/handlers"
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// loadEnvVars is the function responsible to load
// the variables from the .env file.
// It will panic if it can't.
func loadEnvVars() {
	err := godotenv.Load("build/.env")

	if err != nil {
		panic("Couldn't load .env file")
	}
}

// startupDbMigration is the function responsible for doing
// a migration on the database at startup.
// If it isn't able to get access to the database, it will
// panic and print the error.
func startupDbMigration() {
	err := tools.GetDb()
	if err != nil {
		panic(err.Error())
	}

	tools.DB.AutoMigrate(&tools.Role{})
	tools.DB.AutoMigrate(&tools.User{})
}

func main() {
	fmt.Println("Starting themedia API...")

	loadEnvVars()
	startupDbMigration()

	router := gin.Default()

	router.GET("/u/:uuid", handlers.GetUserWithUUID)
	router.POST("/u", handlers.PostUser)
	router.POST("/u/login", handlers.PostLogin)
	router.POST("/u/logout", handlers.PostLogout)

	fmt.Println("API started!")
	router.Run("localhost:8080")
}
