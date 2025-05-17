package middleware

import (
	"net/http"

	ban_tools "github.com/eFournierRobert/themedia/internal/tools/ban"
	"github.com/gin-gonic/gin"
)

// BanCheck is the function that checks if the logged in user
// is banned or not. If yes, it will abort the request with
// an HTTP 401. If not, it will proceeds to the next step of
// the request.
func BanCheck(context *gin.Context) {
	uuid, _ := context.Get("userUUID")

	if ban_tools.IsUserBanned(uuid.(string)) {
		context.AbortWithStatus(http.StatusUnauthorized)
	} else {
		context.Next()
	}
}
