package post_tools

import (
	"errors"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/google/uuid"
)

func GetPostByUUID(uuid *string) (*dbmodels.Post, error) {
	var post dbmodels.Post
	tools.DB.Select("uuid", "title", "body", "user_id").Where("uuid = ?", *uuid).First(&post)
	if post.UUID == "" {
		return nil, errors.New("post was not found")
	}

	return &post, nil
}

func GetAllPost(offset int, limit int) *[]dbmodels.Post {
	if limit == 0 {
		limit = 10
	}
	var postList []dbmodels.Post
	tools.DB.Select("uuid", "title", "body", "user_id").Offset(offset).Limit(limit).Find(postList)

	return &postList
}

func CreatePost(title *string, body *string, userUUID *string, parentPostUUID *string) (*dbmodels.Post, error) {
	var parentPostID *uint = nil

	if parentPostUUID != nil {
		var parentPost dbmodels.Post
		tools.DB.Where("uuid = ?", *parentPostUUID).First(&parentPost)
		if parentPost.ID == 0 {
			return nil, errors.New("parent post not found")
		}

		parentPostID = &parentPost.ID
	}

	var userID uint
	var user dbmodels.User
	tools.DB.Where("uuid = ?", userUUID).First(&user)
	userID = user.ID

	newPost := dbmodels.Post{
		UUID:   uuid.NewString(),
		Title:  title,
		Body:   *body,
		UserID: userID,
		PostID: parentPostID,
	}
	tools.DB.Create(&newPost)

	return &newPost, nil
}

func DeletePost(uuid *string) error {
	var post dbmodels.Post
	tools.DB.Where("uuid = ?", *uuid).First(&post)

	if post.ID == 0 {
		return errors.New("post was not found")
	}

	tools.DB.Unscoped().Delete(&post)
	return nil
}
