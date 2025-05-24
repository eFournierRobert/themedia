package init_tools

import (
	"fmt"
	"os"
	"testing"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	"github.com/eFournierRobert/themedia/internal/tools"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func SetupDatabase(t *testing.T) func(t *testing.T) {
	// Creating temp directory in /tmp
	tempDir, err := os.MkdirTemp("", "themedia-testing-*")
	if err != nil {
		fmt.Println("Couldn't create test directory for SQlite database.")
		os.Exit(1)
	}

	//Creating SQLite database
	db, err := gorm.Open(sqlite.Open(tempDir+"/testing.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Println("Couldn't create test SQLite database.")
		os.Exit(1)
	}

	// Migrate tables and create dummy data
	tableMigration(db)

	CheckIfFirstStartup(db)

	createUsers(db)
	createPosts(db)

	tools.DB = db

	// Delete the directory after test
	return func(t *testing.T) {
		os.RemoveAll(tempDir)
	}
}

func tableMigration(db *gorm.DB) {
	db.AutoMigrate(&dbmodels.Role{})
	db.AutoMigrate(&dbmodels.User{})
	db.AutoMigrate(&dbmodels.Ban{})
	db.AutoMigrate(&dbmodels.Post{})
}

func createUsers(db *gorm.DB) {
	users := []*dbmodels.User{
		{UUID: "de0c8142-5973-478b-9287-37ff25e4e332", Username: "John Doe", PasswordHash: []byte("test"), RoleID: 1, Bio: "Bio of John Doe"},
		{UUID: "35ad671e-0fa0-4829-ae8e-37043d95fc33", Username: "Bright Horizon", PasswordHash: []byte("test"), RoleID: 2, Bio: "Bio of Bright Horizon"},
		{UUID: "dd1614ee-e26f-4949-ba0f-fd8d7df031d2", Username: "Tux Gnu", PasswordHash: []byte("test"), RoleID: 2, Bio: "Bio of Tux Gnu"},
	}
	db.Create(users)
}

func createPosts(db *gorm.DB) {
	testTitles := []string{
		"Test post 1",
		"Test post 2",
		"Test post 3",
	}
	var parentTestPost uint = 1
	posts := []*dbmodels.Post{
		{UUID: "e3631cac-e80d-4908-b902-9e70492079f4", Title: &testTitles[0], Body: "This is the first test post", UserID: 2, PostID: nil},
		{UUID: "8be57d3d-8a55-4bdc-b2e5-e13fe282a467", Title: &testTitles[1], Body: "This is the second test post", UserID: 3, PostID: nil},
		{UUID: "56b757f0-35ec-4055-bf3c-22186a75a3a3", Title: &testTitles[2], Body: "This is the third test post", UserID: 4, PostID: nil},
		{UUID: "a8399ae9-14e6-441b-814c-fe6ce983c8d4", Title: nil, Body: "This is a test answer", UserID: 4, PostID: &parentTestPost},
		{UUID: "1eb075f3-448d-4111-83d9-4f757eea373f", Title: nil, Body: "This is another test answer", UserID: 3, PostID: &parentTestPost},
	}
	db.Create(posts)

}
