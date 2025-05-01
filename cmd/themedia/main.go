package main

import (
	"fmt"
	"net/http"
	
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/eFournierRobert/themedia/internal/handlers"
	"github.com/gin-gonic/gin"
)

func index(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "Hello World")
}

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

	router.GET("/", index)
	router.POST("/users", handlers.PostUser)

	fmt.Println("API started!")
	router.Run("localhost:8080")
}
