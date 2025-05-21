// Package main is the main package of the API.
// The app starts here.
package main

import (
	"fmt"

	user_handlers "github.com/eFournierRobert/themedia/internal/handlers/user"
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

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	user_handlers.AddEndpointsToRouter(router)

	fmt.Println("API started!")
	router.Run("localhost:8080")
}
