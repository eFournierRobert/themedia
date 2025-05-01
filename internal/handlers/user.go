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
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse { 
			Message: "Couldn't create user",
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse { 
			Message: "Populate username and password for creation",
		})
		return
	}

	role, err := tools.FindRole(&user.RoleUUID)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse { 
			Message: "Unknown error",
		})
		return
	}

	createdUser, err := tools.CreateUser(&user.Username, &user.Password, role)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse { 
			Message: "Unknown error",
		})
		return
	}

	context.IndentedJSON(http.StatusCreated, models.UserPostResponse{
		UUID: createdUser.UUID,
		Username: createdUser.Username,
	})
}
