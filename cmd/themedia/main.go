// Package main is the main package of the API.
// The app starts here.
package main

import (
	"fmt"

	handlers "github.com/eFournierRobert/themedia/internal/handlers/user"
	"github.com/eFournierRobert/themedia/internal/middleware"
	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
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

func main() {
	fmt.Println("Starting themedia API...")

	loadEnvVars()
	init_tools.StartupDbMigration()

	router := gin.Default()

	router.GET("/u/:uuid", handlers.GetUserWithUUID)
	router.POST("/u", handlers.PostUser)
	router.POST("/u/login", handlers.PostLogin)
	router.POST("/u/logout", handlers.PostLogout)
	router.DELETE(
		"/u/:uuid",
		middleware.Authorization,
		middleware.BanCheck,
		middleware.AdminOrLoggedInUserCheck,
		handlers.DeleteUser,
	)
	router.PUT(
		"/u/:uuid",
		middleware.Authorization,
		middleware.BanCheck,
		middleware.AdminOrLoggedInUserCheck,
		handlers.PutUser,
	)
	router.POST(
		"/u/:uuid/ban",
		middleware.Authorization,
		middleware.BanCheck,
		middleware.AdminCheck,
		handlers.PostBan,
	)

	fmt.Println("API started!")
	router.Run("localhost:8080")
}
