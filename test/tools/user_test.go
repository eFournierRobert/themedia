package tools

import (
	"testing"

	"github.com/eFournierRobert/themedia/internal/tools"
)

func TestCreateValidUser(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	username := "username"
	password := "password"
	_, err := tools.CreateUser(&username, &password, getAdminRole())
	if err != nil {
		t.Errorf("Couldn't create valid user. Got error %s", err.Error())
	}

	var usr tools.User
	tools.DB.Where("username = ?", username).First(&usr)

	if usr.ID == 0 {
		t.Errorf("Valid user was not created in database %s", err.Error())
	}
}

func getAdminRole() *tools.Role {
	var role tools.Role
	tools.DB.Where("id = ?", 1).First(&role)

	return &role
}
