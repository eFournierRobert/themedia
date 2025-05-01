package handlers

import (
	"net/http"

	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/eFournierRobert/themedia/internal/models"
	"github.com/gin-gonic/gin"
)

func PostUser(context *gin.Context) {
	var user models.UserPost
	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, "Couldn't create user")
		return
	}

	if user.Username == "" || user.Password == "" {
		context.IndentedJSON(http.StatusBadRequest, "Populate username and password for creation")
		return
	}

	createdUser, err := tools.CreateUser(user.Username, user.Password)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, "Unknown error")
		return
	}

	context.IndentedJSON(http.StatusCreated, models.UserPostResponse{
		UUID: createdUser.UUID,
		Username: createdUser.Username,
	})
}
