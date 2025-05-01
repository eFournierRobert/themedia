package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, "Hello World")
}

func main() {
	fmt.Println("Starting themedia API...")

	router := gin.Default()

	router.GET("/", index)


	fmt.Println("API started!")

	router.Run("localhost:8080")
}
