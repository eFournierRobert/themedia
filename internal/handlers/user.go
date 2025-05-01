package handlers

import (
	"net/http"
	"unicode/utf8"

	"github.com/eFournierRobert/themedia/internal/models"
	"github.com/eFournierRobert/themedia/internal/tools"
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

	role, err := tools.FindRoleByUUID(&user.RoleUUID)
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
		Role: models.RoleResponse {
			UUID: role.UUID,
			Name: role.Name,
		},
	})
}

func GetUserWithUUID(context *gin.Context) {
	uuid := context.Param("uuid")
	if utf8.RuneCountInString(uuid) != 36 {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse {
			Message: "Please submit a valid UUID",
		})
	}

	
}
