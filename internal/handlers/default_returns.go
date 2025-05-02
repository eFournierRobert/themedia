// Package handlers is the package that handles the API calls.
// All functions in the package returns the HTTP code and the
// JSON response.
package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/eFournierRobert/themedia/internal/models"
)

// UnknownError is the function responsible to return an internal server error
// with the error message "Unknown error".
func UnknownError(context *gin.Context) {
	context.IndentedJSON(http.StatusInternalServerError, models.ErrorResponse { 
		Message: "Unknown error",
	})
	return
}
