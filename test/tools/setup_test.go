package tools

import (
	"fmt"
	"os"
	"testing"

	"github.com/eFournierRobert/themedia/internal/tools"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var TempDir string

func SetupDatabase(t testing.T) func(t testing.T) {
	// Creating temp directory in /tmp
	tempDir, err := os.MkdirTemp("", "themedia-testing-*")
	if err != nil {
		fmt.Println("Couldn't create test directory for SQlite database.")
		os.Exit(1)
	}

	TempDir = tempDir

	//Creating SQLite database
	db, err := gorm.Open(sqlite.Open(tempDir+"/testing.db"), &gorm.Config{})
	if err != nil {
		fmt.Println("Couldn't create test SQLite database.")
		os.Exit(1)
	}

	// Migrate tables and create dummy data
	db.AutoMigrate(&tools.Role{})
	db.AutoMigrate(&tools.User{})
	db.AutoMigrate(&tools.Ban{})

	tools.CheckIfFirstStartup(db)

	users := []*tools.User{
		{UUID: "de0c8142-5973-478b-9287-37ff25e4e332", Username: "John Doe", PasswordHash: []byte("test"), RoleID: 1, Bio: ""},
		{UUID: "35ad671e-0fa0-4829-ae8e-37043d95fc33", Username: "Bright Horizon", PasswordHash: []byte("test"), RoleID: 2, Bio: ""},
		{UUID: "dd1614ee-e26f-4949-ba0f-fd8d7df031d2", Username: "Tux Gnu", PasswordHash: []byte("test"), RoleID: 2, Bio: ""},
	}
	db.Create(users)

	tools.DB = db

	// Delete the directory after test
	return func(t testing.T) {
		os.RemoveAll(tempDir)
	}
}

func TestSetupDummyDatabase(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	if _, err := os.Stat(TempDir); os.IsNotExist(err) {
		t.Errorf("Directory %s couldn't be created", TempDir)
	}
}
