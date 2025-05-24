package post_tools

import (
	"testing"

	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
	"github.com/stretchr/testify/assert"
)

func TestFindValidPostByUUID(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	UUIDToGet := "8be57d3d-8a55-4bdc-b2e5-e13fe282a467"

	post, err := GetPostByUUID(&UUIDToGet)
	if err != nil {
		t.Errorf("Post was not found. Got %s", err.Error())
	}

	assert.Equal(UUIDToGet, post.UUID)
}

func TestFindPostByInvalidUUID(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	UUIDToGet := "pomme"

	post, err := GetPostByUUID(&UUIDToGet)
	if err == nil {
		t.Errorf("Post with invalid UUID was found")
	}

	assert.Empty(post)
}

func TestFindAllPostWithNoLimitOrOffset(t *testing.T) {
	assert := assert.New(t)
	teardownTest := init_tools.SetupDatabase(t)
	defer teardownTest(t)

	posts := GetAllPost(0, 0)

	assert.Len(posts, 4)
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

	posts := GetAllPost(4, 0)

	assert.Len(posts, 0)
}
