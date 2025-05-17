package user_handlers

import (
	"testing"

	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
	"github.com/gin-gonic/gin"
)

func TestPostUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	router := gin.Default()
	AddEndpointsToRouter(router)

}
