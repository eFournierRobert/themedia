package post_tools

import (
	"testing"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	"github.com/eFournierRobert/themedia/internal/tools"
	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
	"github.com/stretchr/testify/assert"
)

func TestFindValidPostByUUID(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	UUIDToGet := "8be57d3d-8a55-4bdc-b2e5-e13fe282a467"

	post, err := GetPostByUUID(&UUIDToGet)

	assert.NoError(err)
	assert.Equal(UUIDToGet, post.UUID)
}

func TestFindPostByInvalidUUID(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	UUIDToGet := "pomme"

	_, err := GetPostByUUID(&UUIDToGet)

	assert.Error(err)
}

func TestFindAllPostWithNoLimitOrOffset(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	posts := GetAllPost(0, 0)

	assert.Len(posts, 5)
}

func TestFindAllPostWithALimitOfTwoAndAnOffsetOf2(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	posts := GetAllPost(2, 2)

	assert.Len(posts, 2)
}

func TestFindAllPostWithOffsetBiggerThanNumberOfPost(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	posts := GetAllPost(5, 0)

	assert.Len(posts, 0)
}

func TestCreateValidPost(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	title := "new post title"
	body := "I broke my Arch install"
	userUUID := "de0c8142-5973-478b-9287-37ff25e4e332"
	post, err := CreatePost(&title, &body, &userUUID, nil)

	assert.NoError(err)
	assert.Equal(title, *post.Title)
	assert.Equal(body, post.Body)
}

func TestCreatePostWithInvalidUser(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	title := "new post title"
	body := "I broke my Arch install"
	userUUID := "womp womp"
	_, err := CreatePost(&title, &body, &userUUID, nil)

	assert.Error(err)
}

func TestCreateAnswerPost(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	body := "This is the greatest post ever made"
	userUUID := "de0c8142-5973-478b-9287-37ff25e4e332"
	parentPostUUID := "e3631cac-e80d-4908-b902-9e70492079f4"
	post, err := CreatePost(nil, &body, &userUUID, &parentPostUUID)

	assert.NoError(err)
	assert.Nil(post.Title)
	assert.Equal(body, post.Body)
	assert.NotEmpty(post.PostID)
}

func TestCreateAnswerPostWithInvalidParentPostUUID(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	body := "This is the greatest post ever made"
	userUUID := "de0c8142-5973-478b-9287-37ff25e4e332"
	parentPostUUID := "Bestest of best post"
	_, err := CreatePost(nil, &body, &userUUID, &parentPostUUID)

	assert.Error(err)
}

func TestDeleteValidPost(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	UUID := "e3631cac-e80d-4908-b902-9e70492079f4"
	DeletePost(&UUID)

	var post dbmodels.Post
	tools.DB.Where("uuid = ?", UUID).First(&post)

	assert.Empty(post)

	// Check if answer post was deleted too
	tools.DB.Where("uuid = ?", "a8399ae9-14e6-441b-814c-fe6ce983c8d4").First(&post)
	assert.Empty(post)
}

func TestDeleteInvalidPost(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	UUID := "coconut"
	err := DeletePost(&UUID)

	assert.Error(err)
}

func TestGetValidPostThread(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	parentUUID := "e3631cac-e80d-4908-b902-9e70492079f4"
	thread, err := GetPostThread(&parentUUID, 0, 0)

	assert.NoError(err)
	assert.Len(thread, 2)
	assert.Equal("a8399ae9-14e6-441b-814c-fe6ce983c8d4", thread[0].UUID)
	assert.Equal("1eb075f3-448d-4111-83d9-4f757eea373f", thread[1].UUID)
}

func TestGetPostThreadFromParentPostWithoutAnswers(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	parentUUID := "8be57d3d-8a55-4bdc-b2e5-e13fe282a467"
	thread, err := GetPostThread(&parentUUID, 0, 0)

	assert.NoError(err)
	assert.Len(thread, 0)
}

func TestGetPostThreadFromInvalidParentPost(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	parentUUID := "Fraise"
	thread, err := GetPostThread(&parentUUID, 0, 0)

	assert.Error(err)
	assert.Nil(thread)
}

func TestGetPostFromValidUser(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	userUUID := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	posts := GetAllPostFromUser(&userUUID, 0, 0)

	assert.Len(posts, 2)
}

func TestGetPostFromInvalidUser(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	userUUID := "Poire"
	posts := GetAllPostFromUser(&userUUID, 0, 0)

	assert.Len(posts, 0)
}
