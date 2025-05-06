// Package handlers is the package that handles the API calls.
// All functions in the package returns the HTTP code and the
// JSON response.
package handlers

import (
	"net/http"
	"os"
	"time"
	"unicode/utf8"

	"github.com/eFournierRobert/themedia/internal/models"
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// PostUser is the function that handles the API POST /u.
// It will create a new user in the database.
func PostUser(context *gin.Context) {
	var user models.UserPost
	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Couldn't create user",
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
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
		UUID:     createdUser.UUID,
		Username: createdUser.Username,
		Role: models.RoleResponse{
			UUID: role.UUID,
			Name: role.Name,
		},
	})
}

// GetUserWithUUID is the function that handles the API GET /u/{uuid}.
// It will find the user with that UUID in the database with its role.
func GetUserWithUUID(context *gin.Context) {
	uuid := context.Param("uuid")
	if !validUUIDCheck(&uuid) {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Please submit a valid UUID",
		})
		return
	}

	fullUser, err := tools.FindFullUserByUUID(&uuid)
	if err != nil {
		UnknownError(context)
		return
	}

	if fullUser == nil || fullUser.RoleUUID == "" {
		context.IndentedJSON(http.StatusNotFound, models.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	context.IndentedJSON(http.StatusFound, models.UserResponse{
		UUID:     fullUser.UserUUID,
		Username: fullUser.Username,
		Role: models.RoleResponse{
			UUID: fullUser.RoleUUID,
			Name: fullUser.Name,
		},
	})
}

// PostLogin is the function that handles the API POST /u/login.
// It will check the credentials and create a cookie with the
// JWT token inside.
func PostLogin(context *gin.Context) {
	var user models.UserPost
	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Couldn't read credentials",
		})
		return
	}

	if tools.IsUserBanned(user.UUID) {
		context.IndentedJSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "User temporarily banned",
		})
		return
	}

	isCorrect, err := tools.VerifyPassword(&user.UUID, &user.Password)
	if err != nil || !isCorrect {
		context.IndentedJSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Login failed",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UUID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		context.IndentedJSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Login failed",
		})
		return
	}

	context.SetSameSite(http.SameSiteLaxMode)
	context.SetCookie("Authorization", tokenString, 3600*12, "", "", true, true)

	context.IndentedJSON(http.StatusOK, "Login successful")
}

// PostLogout is the function that logs out the user.
// It sets the maximum age of the cookie containing the
// JWT token to -1.
func PostLogout(context *gin.Context) {
	context.SetCookie("Authorization", "", -1, "", "", true, true)
	context.IndentedJSON(http.StatusOK, "Logout successful")
}

// DeleteUser is the function that deletes a given user.
func DeleteUser(context *gin.Context) {
	uuid := context.Param("uuid")
	if !validUUIDCheck(&uuid) {
		context.IndentedJSON(http.StatusBadRequest, models.ErrorResponse{
			Message: "Please submit a valid UUID",
		})
		return
	}

	err := tools.DeleteUser(uuid)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, models.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, "User deleted")
}

// PutUser is the function that updates a given user.
func PutUser(context *gin.Context) {
	var user models.UserPost
	context.BindJSON(&user)

	loggedUserUUID, _ := context.Get("userUUID")

	if user.RoleUUID != "" && !tools.IsUserAdmin(loggedUserUUID.(string)) {
		context.IndentedJSON(http.StatusForbidden, models.ErrorResponse{
			Message: "Need to be admin to update role.",
		})
		return
	}

	err := tools.UpdateUser(context.Param("uuid"), &user)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "Unknown error",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, "User updated")
}

// PostBan is the function that handles temp banning of the
// user in the request.
func PostBan(context *gin.Context) {
	var ban models.Ban
	context.BindJSON(&ban)

	err := tools.CreateBan(context.Param("uuid"), ban.EndDatetime)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: err.Error(),
		})
	}

	context.IndentedJSON(http.StatusOK, "User temporarily banned")
}

// validUUIDCheck is a function that returns true if
// a UUID is valid and false if it isn't.
func validUUIDCheck(uuid *string) bool {
	return utf8.RuneCountInString(*uuid) == 36
}
