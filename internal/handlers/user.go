// Package handlers is the package that handles the API calls. 
// All functions in the package returns the HTTP code and the
// JSON response.
package handlers

import (
	"net/http"
	"unicode/utf8"

	"github.com/eFournierRobert/themedia/internal/models"
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/gin-gonic/gin"
)

// PostUser is the function that handles the API POST /u.
// It will create a new user in the database.
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
		UnknownError(context)
		return
	}

	createdUser, err := tools.CreateUser(&user.Username, &user.Password, role)
	if err != nil {
		UnknownError(context)
		return
	}

	context.IndentedJSON(http.StatusCreated, models.UserResponse{
		UUID: createdUser.UUID,
		Username: createdUser.Username,
		Role: models.RoleResponse {
			UUID: role.UUID,
			Name: role.Name,
		},
	})
}

// GetUserWithUUID is the function that handles the API GET /u/{uuid}.
// It will find the user with that UUID in the database with its role.
func GetUserWithUUID(context *gin.Context) {
	uuid := context.Param("uuid")
	if utf8.RuneCountInString(uuid) != 36 {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse {
			Message: "Please submit a valid UUID",
		})
	}

	fullUser, err := tools.FindFullUserByUUID(&uuid)
	if err != nil {
		UnknownError(context)
		return
	}

	if fullUser == nil{
		context.IndentedJSON(http.StatusNotFound, models.ErrorResponse {
			Message: "User not found",
		})
		return
	}
	
	context.IndentedJSON(http.StatusFound, models.UserResponse {
		UUID: fullUser.UserUUID,
		Username: fullUser.Username,
		Role: models.RoleResponse{
			UUID: fullUser.RoleUUID,
			Name: fullUser.Name,
		},
	})
}
