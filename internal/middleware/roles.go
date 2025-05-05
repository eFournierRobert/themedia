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
