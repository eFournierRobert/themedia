// Package main is the main package of the API.
// The app starts here.
package main

import (
	"fmt"
	
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/eFournierRobert/themedia/internal/handlers"
	"github.com/gin-gonic/gin"
)

// startupDbMigration is the function responsible for doing
// a migration on the database at startup.
// If it isn't able to get access to the database, it will
// panic and print the error.
func startupDbMigration() {
	db, err := tools.GetDb()
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&tools.Role{})
	db.AutoMigrate(&tools.User{})
}

func main() {
	fmt.Println("Starting themedia API...")
	
	startupDbMigration()

	router := gin.Default()

	router.GET("/u/:uuid", handlers.GetUserWithUUID)
	router.POST("/u", handlers.PostUser)

	fmt.Println("API started!")
	router.Run("localhost:8080")
}
