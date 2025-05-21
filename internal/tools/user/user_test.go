package user_tools

import (
	"bytes"
	"testing"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	jsonmodels "github.com/eFournierRobert/themedia/internal/models/json"
	"github.com/eFournierRobert/themedia/internal/tools"
	init_tools "github.com/eFournierRobert/themedia/internal/tools/init"
)

func TestCreateValidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	username := "username"
	password := "password"
	_, err := CreateUser(&username, &password, getAdminRole())
	if err != nil {
		t.Errorf("Couldn't create valid user. Got error %s", err.Error())
	}

	var usr dbmodels.User
	tools.DB.Where("username = ?", username).First(&usr)

	if usr.ID == 0 {
		t.Errorf("Valid user was not created in database %s", err.Error())
	}
}

func TestCreateUserWithUsernameDeleted(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	username := "deleted"
	password := "password"
	_, err := CreateUser(&username, &password, getAdminRole())
	if err == nil {
		t.Errorf("Could create new user with reserved username deleted")
	}
}

func TestCreateUserWithEmptyRole(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	username := "username"
	password := "password"
	_, err := CreateUser(&username, &password, &dbmodels.Role{})
	if err != nil {
		t.Errorf("Couldn't create user with empty role. Got error %s", err.Error())
	}

	var usr dbmodels.User
	tools.DB.Where("username = ?", username).First(&usr)

	if usr.RoleID != 2 {
		t.Errorf("User with empty role was not created with the user role. Got created with, %d", usr.RoleID)
	}
}

func TestFindFullValidUserByUUID(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"
	usr, err := FindFullUserByUUID(&uuid)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if usr.ID == 0 {
		t.Errorf("Couldn't find valid full user with UUID %s", uuid)
	} else if usr.RoleUUID != getAdminRole().UUID {
		t.Errorf("Role UUID doesn't corresponds. Got %s and admin is %s", usr.RoleUUID, getAdminRole().UUID)
	}
}

func TestFindFullUserThatDoesNotExist(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "banane"
	usr, err := FindFullUserByUUID(&uuid)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if usr.ID != 0 {
		t.Errorf("Found user while we weren't supposed to. Checked for UUID %s and got %+v", uuid, usr)
	}
}

func TestFindRoleByName(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	adminName := "admin"
	role, err := FindRoleByName(&adminName)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if role.ID == 0 {
		t.Errorf("Did not find admin role with name %s", adminName)
	}
}

func TestFindInvalidRoleWithName(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "pomme"
	role, err := FindRoleByName(&uuid)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if role.ID != 0 {
		t.Errorf("Got invalid role while searching by name. Searched name %s, got %+v", uuid, role)
	}
}

func TestVerifyValidPasswordForValidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	// Create new user so we have a password hash
	username := "username"
	password := "password"
	usr, err := CreateUser(&username, &password, getAdminRole())
	if err != nil {
		t.Errorf("Couldn't create valid user. Got error %s", err.Error())
	}

	b, err := VerifyPassword(&usr.UUID, &password)
	if err != nil {
		t.Errorf("User was not found with UUID %s", usr.UUID)
	}

	if !b {
		t.Errorf("Good password does not corresponds to password in the database")
	}
}

func TestVerifyInvalidPasswordForValidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	// Create new user so we have a password hash
	username := "username"
	password := "password"
	usr, err := CreateUser(&username, &password, getAdminRole())
	if err != nil {
		t.Errorf("Couldn't create valid user. Got error %s", err.Error())
	}

	wrongPassword := "Chocolatine"
	b, err := VerifyPassword(&usr.UUID, &wrongPassword)
	if err != nil {
		t.Errorf("User was not found with UUID %s", usr.UUID)
	}

	if b {
		t.Errorf("Wrong password worked with test user in the database")
	}
}

func TestVerifyPasswordForInvalidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	invalidUUID := "Chocolatine"
	dummyPassword := "Poutine"
	_, err := VerifyPassword(&invalidUUID, &dummyPassword)
	if err == nil {
		t.Errorf("Invalid user was found with the UUID %s", invalidUUID)
	}
}

func TestDoesValidUserExist(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"

	b := DoesUserExist(uuid)
	if !b {
		t.Errorf("Valid user with UUID %s was not found in the database", uuid)
	}
}

func TestDoesInvalidUserExist(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "Poire"

	b := DoesUserExist(uuid)
	if b {
		t.Errorf("Invalid user was found in the database with the UUID %s", uuid)
	}
}

func TestIsValidUserAdmin(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"

	b := IsUserAdmin(uuid)
	if !b {
		t.Errorf("Valid admin user with UUID %s is not an admin the database", uuid)
	}
}

func TestIsInvalidUserAdmin(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"

	b := IsUserAdmin(uuid)
	if b {
		t.Errorf("User with the role user and the UUID %s was considered an admin", uuid)
	}
}

func TestDeleteValidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	err := DeleteUser(uuid)
	if err != nil {
		t.Errorf("User with UUID %s that was supposed to exist in the database doesn't", uuid)
	}

	if DoesUserExist(uuid) {
		t.Errorf("User with UUID %s is supposed to be deleted from the database but wasn't", uuid)
	}
}

func TestDeleteInvalidUser(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "Pouding chomeur"
	err := DeleteUser(uuid)
	if err == nil {
		t.Errorf("Invalid user with UUID %s was found in the database", uuid)
	}
}

func TestUpdateValidUserWithOnlyUsername(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := jsonmodels.UserPost{
		UUID:     uuid,
		Username: "New Username",
	}

	err := UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if updatedUser.Username != user.Username {
		t.Errorf("Username was not updated. Expected %s, got %s", user.Username, updatedUser.Username)
	}
}

func TestUpdateValidUserWithOnlyPassword(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := jsonmodels.UserPost{
		UUID:     uuid,
		Password: "New password",
	}

	var oldUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&oldUser)

	err := UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if bytes.Equal(updatedUser.PasswordHash, oldUser.PasswordHash) {
		t.Errorf("Password was not updated.")
	}
}

func TestUpdateValidUserWithOnlyBio(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := jsonmodels.UserPost{
		UUID: uuid,
		Bio:  "New bio",
	}

	var oldUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&oldUser)

	err := UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if updatedUser.Bio == oldUser.Bio {
		t.Errorf("Bio was not updated. Expected %s, got %s", user.Bio, updatedUser.Bio)
	}
}

func TestUpdateValidUserWithValidRole(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := jsonmodels.UserPost{
		UUID: uuid,
		Role: "admin",
	}

	var oldUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&oldUser)

	err := UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser dbmodels.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if updatedUser.RoleID == oldUser.RoleID {
		t.Errorf("Role was not updated. Expected %d, got %d", getAdminRole().ID, updatedUser.RoleID)
	}
}

func TestUpdateValidUSerWithInvalidRole(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := jsonmodels.UserPost{
		UUID: uuid,
		Role: "Poire",
	}

	err := UpdateUser(uuid, &user)
	if err == nil {
		t.Errorf("User was updated with invalid role of name %s", user.Role)
	}
}

func TestUpdateUserWithUsernameDeleted(t *testing.T) {
	teardownSuite := init_tools.SetupDatabase(t)
	defer teardownSuite(t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := jsonmodels.UserPost{
		UUID:     uuid,
		Username: "deleted",
	}

	err := UpdateUser(uuid, &user)
	if err == nil {
		t.Errorf("User was updated with reserved username deleted")
	}
}

func getAdminRole() *dbmodels.Role {
	var role dbmodels.Role
	tools.DB.Where("id = ?", 1).First(&role)

	return &role
}
