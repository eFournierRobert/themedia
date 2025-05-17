// Package handlers is the package that handles the API calls.
// All functions in the package returns the HTTP code and the
// JSON response.
package user_handlers

import (
	"net/http"
	"os"
	"time"
	"unicode/utf8"

	"github.com/eFournierRobert/themedia/internal/handlers"
	"github.com/eFournierRobert/themedia/internal/middleware"
	jsonmodels "github.com/eFournierRobert/themedia/internal/models/json"
	ban_tools "github.com/eFournierRobert/themedia/internal/tools/ban"
	user_tools "github.com/eFournierRobert/themedia/internal/tools/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AddEndpointsToRouter(router *gin.Engine) {
	router.GET("/u/:uuid", GetUserWithUUID)
	router.POST("/u", PostUser)
	router.POST("/u/login", PostLogin)
	router.POST("/u/logout", PostLogout)
	router.DELETE(
		"/u/:uuid",
		middleware.Authorization,
		middleware.BanCheck,
		middleware.AdminOrLoggedInUserCheck,
		DeleteUser,
	)
	router.PUT(
		"/u/:uuid",
		middleware.Authorization,
		middleware.BanCheck,
		middleware.AdminOrLoggedInUserCheck,
		PutUser,
	)
	router.POST(
		"/u/:uuid/ban",
		middleware.Authorization,
		middleware.BanCheck,
		middleware.AdminCheck,
		PostBan,
	)
}

// PostUser is the function that handles the API POST /u.
// It will create a new user in the database.
func PostUser(context *gin.Context) {
	var user jsonmodels.UserPost
	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: "Couldn't create user",
		})
		return
	}

	if user.Username == "" || user.Password == "" {
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: "Populate username and password for creation",
		})
		return
	}

	role, err := user_tools.FindRoleByName(&user.Role)
	if err != nil {
		handlers.UnknownError(context)
		return
	}

	createdUser, err := user_tools.CreateUser(&user.Username, &user.Password, role)
	if err != nil {
		handlers.UnknownError(context)
		return
	}

	context.IndentedJSON(http.StatusCreated, jsonmodels.UserResponse{
		UUID:     createdUser.UUID,
		Username: createdUser.Username,
		Role: jsonmodels.RoleResponse{
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
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: "Please submit a valid UUID",
		})
		return
	}

	fullUser, err := user_tools.FindFullUserByUUID(&uuid)
	if err != nil {
		handlers.UnknownError(context)
		return
	}

	if fullUser == nil || fullUser.RoleUUID == "" {
		context.IndentedJSON(http.StatusNotFound, jsonmodels.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	context.IndentedJSON(http.StatusFound, jsonmodels.UserResponse{
		UUID:     fullUser.UserUUID,
		Username: fullUser.Username,
		Role: jsonmodels.RoleResponse{
			UUID: fullUser.RoleUUID,
			Name: fullUser.Name,
		},
	})
}

// PostLogin is the function that handles the API POST /u/login.
// It will check the credentials and create a cookie with the
// JWT token inside.
func PostLogin(context *gin.Context) {
	var user jsonmodels.UserPost
	if err := context.BindJSON(&user); err != nil {
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: "Couldn't read credentials",
		})
		return
	}

	if ban_tools.IsUserBanned(user.UUID) {
		context.IndentedJSON(http.StatusUnauthorized, jsonmodels.ErrorResponse{
			Message: "User temporarily banned",
		})
		return
	}

	isCorrect, err := user_tools.VerifyPassword(&user.UUID, &user.Password)
	if err != nil || !isCorrect {
		context.IndentedJSON(http.StatusUnauthorized, jsonmodels.ErrorResponse{
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
		context.IndentedJSON(http.StatusUnauthorized, jsonmodels.ErrorResponse{
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
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: "Please submit a valid UUID",
		})
		return
	}

	err := user_tools.DeleteUser(uuid)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, jsonmodels.ErrorResponse{
			Message: "User not found",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, "User deleted")
}

// PutUser is the function that updates a given user.
func PutUser(context *gin.Context) {
	var user jsonmodels.UserPost
	context.BindJSON(&user)

	loggedUserUUID, _ := context.Get("userUUID")

	if user.Role != "" && !user_tools.IsUserAdmin(loggedUserUUID.(string)) {
		context.IndentedJSON(http.StatusForbidden, jsonmodels.ErrorResponse{
			Message: "Need to be admin to update role.",
		})
		return
	}

	err := user_tools.UpdateUser(context.Param("uuid"), &user)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, jsonmodels.ErrorResponse{
			Message: "Unknown error",
		})
		return
	}

	context.IndentedJSON(http.StatusOK, "User updated")
}

// PostBan is the function that handles temp banning of the
// user in the request.
func PostBan(context *gin.Context) {
	var ban jsonmodels.Ban
	context.BindJSON(&ban)

	err := ban_tools.CreateBan(context.Param("uuid"), ban.EndDatetime)
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, jsonmodels.ErrorResponse{
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
