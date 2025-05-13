package tools

import (
	"bytes"
	"testing"

	"github.com/eFournierRobert/themedia/internal/models"
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

func TestCreateUserWithUsernameDeleted(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	username := "deleted"
	password := "password"
	_, err := tools.CreateUser(&username, &password, getAdminRole())
	if err == nil {
		t.Errorf("Could create new user with reserved username deleted")
	}
}

func TestCreateUserWithEmptyRole(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	username := "username"
	password := "password"
	_, err := tools.CreateUser(&username, &password, &tools.Role{})
	if err != nil {
		t.Errorf("Couldn't create user with empty role. Got error %s", err.Error())
	}

	var usr tools.User
	tools.DB.Where("username = ?", username).First(&usr)

	if usr.RoleID != 2 {
		t.Errorf("User with empty role was not created with the user role. Got created with, %d", usr.RoleID)
	}
}

func TestFindFullValidUserByUUID(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"
	usr, err := tools.FindFullUserByUUID(&uuid)
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
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "banane"
	usr, err := tools.FindFullUserByUUID(&uuid)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if usr.ID != 0 {
		t.Errorf("Found user while we weren't supposed to. Checked for UUID %s and got %+v", uuid, usr)
	}
}

func TestFindRoleByUUID(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := getAdminRole().UUID
	role, err := tools.FindRoleByUUID(&uuid)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if role.ID == 0 {
		t.Errorf("Did not find admin role with UUID %s", uuid)
	}
}

func TestFindInvalidRoleWithUUID(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "pomme"
	role, err := tools.FindRoleByUUID(&uuid)
	if err != nil {
		t.Errorf("Unknown error while finding full user. Got error %s", err.Error())
	}

	if role.ID != 0 {
		t.Errorf("Got invalid role while searching by UUID. Searched UUID %s, got %+v", uuid, role)
	}
}

func TestVerifyValidPasswordForValidUser(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	// Create new user so we have a password hash
	username := "username"
	password := "password"
	usr, err := tools.CreateUser(&username, &password, getAdminRole())
	if err != nil {
		t.Errorf("Couldn't create valid user. Got error %s", err.Error())
	}

	b, err := tools.VerifyPassword(&usr.UUID, &password)
	if err != nil {
		t.Errorf("User was not found with UUID %s", usr.UUID)
	}

	if !b {
		t.Errorf("Good password does not corresponds to password in the database")
	}
}

func TestVerifyInvalidPasswordForValidUser(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	// Create new user so we have a password hash
	username := "username"
	password := "password"
	usr, err := tools.CreateUser(&username, &password, getAdminRole())
	if err != nil {
		t.Errorf("Couldn't create valid user. Got error %s", err.Error())
	}

	wrongPassword := "Chocolatine"
	b, err := tools.VerifyPassword(&usr.UUID, &wrongPassword)
	if err != nil {
		t.Errorf("User was not found with UUID %s", usr.UUID)
	}

	if b {
		t.Errorf("Wrong password worked with test user in the database")
	}
}

func TestVerifyPasswordForInvalidUser(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	invalidUUID := "Chocolatine"
	dummyPassword := "Poutine"
	_, err := tools.VerifyPassword(&invalidUUID, &dummyPassword)
	if err == nil {
		t.Errorf("Invalid user was found with the UUID %s", invalidUUID)
	}
}

func TestDoesValidUserExist(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"

	b := tools.DoesUserExist(uuid)
	if !b {
		t.Errorf("Valid user with UUID %s was not found in the database", uuid)
	}
}

func TestDoesInvalidUserExist(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "Poire"

	b := tools.DoesUserExist(uuid)
	if b {
		t.Errorf("Invalid user was found in the database with the UUID %s", uuid)
	}
}

func TestIsValidUserAdmin(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "de0c8142-5973-478b-9287-37ff25e4e332"

	b := tools.IsUserAdmin(uuid)
	if !b {
		t.Errorf("Valid admin user with UUID %s is not an admin the database", uuid)
	}
}

func TestIsInvalidUserAdmin(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"

	b := tools.IsUserAdmin(uuid)
	if b {
		t.Errorf("User with the role user and the UUID %s was considered an admin", uuid)
	}
}

func TestDeleteValidUser(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	err := tools.DeleteUser(uuid)
	if err != nil {
		t.Errorf("User with UUID %s that was supposed to exist in the database doesn't", uuid)
	}

	if tools.DoesUserExist(uuid) {
		t.Errorf("User with UUID %s is supposed to be deleted from the database but wasn't", uuid)
	}
}

func TestDeleteInvalidUser(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "Pouding chomeur"
	err := tools.DeleteUser(uuid)
	if err == nil {
		t.Errorf("Invalid user with UUID %s was found in the database", uuid)
	}
}

func TestUpdateValidUserWithOnlyUsername(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := models.UserPost{
		UUID:     uuid,
		Username: "New Username",
	}

	err := tools.UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if updatedUser.Username != user.Username {
		t.Errorf("Username was not updated. Expected %s, got %s", user.Username, updatedUser.Username)
	}
}

func TestUpdateValidUserWithOnlyPassword(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := models.UserPost{
		UUID:     uuid,
		Password: "New password",
	}

	var oldUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&oldUser)

	err := tools.UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if bytes.Equal(updatedUser.PasswordHash, oldUser.PasswordHash) {
		t.Errorf("Password was not updated.")
	}
}

func TestUpdateValidUserWithOnlyBio(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := models.UserPost{
		UUID: uuid,
		Bio:  "New bio",
	}

	var oldUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&oldUser)

	err := tools.UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if updatedUser.Bio == oldUser.Bio {
		t.Errorf("Bio was not updated. Expected %s, got %s", user.Bio, updatedUser.Bio)
	}
}

func TestUpdateValidUserWithValidRole(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := models.UserPost{
		UUID:     uuid,
		RoleUUID: getAdminRole().UUID,
	}

	var oldUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&oldUser)

	err := tools.UpdateUser(uuid, &user)
	if err != nil {
		t.Errorf("Error during user update. Got %s", err.Error())
	}

	var updatedUser tools.User
	tools.DB.Where("uuid = ?", uuid).First(&updatedUser)

	if updatedUser.ID == 0 {
		t.Errorf("Updated user was not found using valid UUID %s", uuid)
	} else if updatedUser.RoleID == oldUser.RoleID {
		t.Errorf("Role was not updated. Expected %d, got %d", getAdminRole().ID, updatedUser.RoleID)
	}
}

func TestUpdateValidUSerWithInvalidRole(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := models.UserPost{
		UUID:     uuid,
		RoleUUID: "Poire",
	}

	err := tools.UpdateUser(uuid, &user)
	if err == nil {
		t.Errorf("User was updated with invalid role of UUID %s", user.RoleUUID)
	}
}

func TestUpdateUserWithUsernameDeleted(t *testing.T) {
	teardownSuite := SetupDatabase(*t)
	defer teardownSuite(*t)

	uuid := "35ad671e-0fa0-4829-ae8e-37043d95fc33"
	user := models.UserPost{
		UUID:     uuid,
		Username: "deleted",
	}

	err := tools.UpdateUser(uuid, &user)
	if err == nil {
		t.Errorf("User was updated with reserved username deleted")
	}
}

func getAdminRole() *tools.Role {
	var role tools.Role
	tools.DB.Where("id = ?", 1).First(&role)

	return &role
}
