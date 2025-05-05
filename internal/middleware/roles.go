package middleware

import (
	"net/http"

	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/gin-gonic/gin"
)

// AdminCheck checks if the connected user is an admin or not.
// If yes, it will continue the request, if not, it will abort and return
// HTTP 401.
func AdminCheck(context *gin.Context) {
	userUUID, exist := context.Get("userUUID")
	if exist == false {
		context.AbortWithStatus(http.StatusUnauthorized)
	}

	if tools.IsUserAdmin(userUUID.(string)) {
		context.Next()
	} else {
		context.AbortWithStatus(http.StatusUnauthorized)
	}
}

// AdminOrLoggedInUserCheck checks if the connected user is an
// admin or if the UUID is the request corresponds to the logged in user.
// If not, it will abort and reutrn HTTP 401.
func AdminOrLoggedInUserCheck(context *gin.Context) {
	userUUID, exist := context.Get("userUUID")
	if exist == false {
		context.AbortWithStatus(http.StatusUnauthorized)
	}

	if userUUID.(string) == context.Param("uuid") || tools.IsUserAdmin(userUUID.(string)) {
		context.Next()
	} else {
		context.AbortWithStatus(http.StatusUnauthorized)
	}
}
