package main

import (
	"fmt"
	
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/eFournierRobert/themedia/internal/handlers"
	"github.com/gin-gonic/gin"
)

func startupDbMigration() {
	db, err := tools.GetDb()
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&tools.User{})
}

func main() {
	fmt.Println("Starting themedia API...")
	
	startupDbMigration()

	router := gin.Default()

	router.POST("/u", handlers.PostUser)

	fmt.Println("API started!")
	router.Run("localhost:8080")
}
