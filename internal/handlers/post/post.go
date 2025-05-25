package post_handlers

import (
	"net/http"
	"strconv"
	"unicode/utf8"

	jsonmodels "github.com/eFournierRobert/themedia/internal/models/json"
	post_tools "github.com/eFournierRobert/themedia/internal/tools/post"
	user_tools "github.com/eFournierRobert/themedia/internal/tools/user"
	"github.com/gin-gonic/gin"
)

func AddEndpointsToRouter(router *gin.Engine) {
	router.GET("/p/:uuid", GetPostWithUUID)
	router.GET("/p", GetAllPost)
}

func GetPostWithUUID(context *gin.Context) {
	uuid := context.Param("uuid")
	if !validUUIDCheck(&uuid) {
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: "Please submit a valid UUID",
		})
		return
	}

	post, err := post_tools.GetPostByUUID(&uuid)
	if err != nil {
		context.IndentedJSON(http.StatusBadRequest, jsonmodels.ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userUUID, _ := user_tools.FindUserUUIDWithID(post.UserID)

	context.IndentedJSON(http.StatusFound, jsonmodels.Post{
		UUID:     post.UUID,
		Title:    *post.Title,
		Body:     post.Body,
		UserUUID: *userUUID,
	})

}

func GetAllPost(context *gin.Context) {
	offset, limit := parseOffsetAndLimit(context)

	posts := post_tools.GetAllPost(offset, limit)

	var jsonPosts []jsonmodels.Post
	for _, post := range posts {
		userUUID, _ := user_tools.FindUserUUIDWithID(post.UserID)
		jsonPosts = append(jsonPosts, jsonmodels.Post{
			UUID:     post.UUID,
			Title:    *post.Title,
			Body:     post.Body,
			UserUUID: *userUUID,
		})
	}

	context.IndentedJSON(http.StatusOK, jsonPosts)
}

func parseOffsetAndLimit(context *gin.Context) (int, int) {
	paramPairs := context.Request.URL.Query()
	offsetString := paramPairs.Get("offset")
	limitString := paramPairs.Get("limit")

	var offset int
	var err error
	var limit int

	if offsetString != "" {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			offset = 0
		}
	} else {
		offset = 0
	}

	if limitString != "" {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			limit = 0
		}
	} else {
		limit = 0
	}

	return offset, limit
}

// validUUIDCheck is a function that returns true if
// a UUID is valid and false if it isn't.
func validUUIDCheck(uuid *string) bool {
	return utf8.RuneCountInString(*uuid) == 36
}
