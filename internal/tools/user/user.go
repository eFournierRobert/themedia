// Package tools is the package containing
// all the request made to the database.
package user_tools

import (
	"errors"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	jsonmodels "github.com/eFournierRobert/themedia/internal/models/json"
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser is the function responsible for inserting a new user in the database.
// the password passed must be in plain text since function will be hashing.
// It returns a pointer to a User struct or an error.
func CreateUser(username *string, password *string, role *dbmodels.Role) (*dbmodels.User, error) {
	if *username == "deleted" {
		return nil, errors.New("cannot create a user with the username deleted")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("couldn't hash password")
	}

	if role.ID == 0 {
		tools.DB.Where("name = ?", "user").First(&role)
		if role.ID == 0 {
			return nil, errors.New("couldn't find the user role")
		}
	}

	user := dbmodels.User{
		UUID:         uuid.NewString(),
		Username:     *username,
		PasswordHash: hashedPassword,
		RoleID:       role.ID,
	}

	tools.DB.Create(&user)

	return &user, nil
}

// FindFullUserByUUID is the function responsible for finding the user that
// has the given UUID in the database.
// It will return a pointer to a FullUser struct or an error.
func FindFullUserByUUID(uuid *string) (*dbmodels.FullUser, error) {
	var fullUser dbmodels.FullUser
	tools.DB.Table("users").Select(
		"users.id",
		"users.uuid AS user_uuid",
		"users.username",
		"users.bio",
		"roles.uuid AS role_uuid",
		"roles.name",
	).Where("users.uuid = ?", *uuid).Joins("JOIN roles ON roles.id = users.role_id").First(&fullUser)

	return &fullUser, nil
}

// FindRoleByUUID is the function responsible for finding
// the role that has the given UUID in the database.
// It will return a pointer to a Role struct or an error.
func FindRoleByUUID(uuid *string) (*dbmodels.Role, error) {
	var role dbmodels.Role
	tools.DB.Where("uuid = ?", *uuid).First(&role)

	return &role, nil
}

// VerifyPassword is the function that receives a raw passwords and
// compares it to the hashed one in the database for the user with the given
// UUID. It will return a boolean or an error.
func VerifyPassword(uuid *string, password *string) (bool, error) {
	var user dbmodels.User
	tools.DB.Table("users").Select("id", "password_hash").Where("uuid = ?", uuid).First(&user)

	if user.ID == 0 {
		return false, errors.New("did not find user")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*password))
	if err != nil {
		return false, nil
	}

	return true, nil
}

// DoesUserExist is the function that checks if a user with the given
// UUID exist in the database. It returns a boolean.
func DoesUserExist(uuid string) bool {
	var user dbmodels.User
	tools.DB.Table("users").Select("id").Where("uuid = ?", uuid).First(&user)

	return user.ID != 0
}

// IsUserAdmin takes the given user UUID and checks if the user
// is an admin or not.
func IsUserAdmin(uuid string) bool {
	var user dbmodels.FullUser
	tools.DB.Table("users").Select(
		"users.id",
		"roles.name",
	).Where("users.uuid = ?", uuid).Joins("JOIN roles on roles.id = users.role_id").First(&user)

	return user.ID != 0 && user.Name == "admin"
}

// DeleteUser takes the given user UUID and hard delete it
// from the database.
func DeleteUser(uuid string) error {
	var user dbmodels.User
	tools.DB.Where("uuid = ? AND username != 'deleted'", uuid).First(&user)

	if user.ID == 0 {
		return errors.New("User does not exist")
	}

	tools.DB.Unscoped().Delete(&user)
	return nil
}

// UpdateUser updates the user in the database with the new information received.
// Returns nil if successful or an error if not.
func UpdateUser(uuid string, user *jsonmodels.UserPost) error {
	if user.Username == "deleted" {
		return errors.New("cannot modify username for deleted")
	}

	var oldUser dbmodels.User
	var updatedUser dbmodels.User
	tools.DB.Table("users").Select("id").Where("uuid = ?", uuid).First(&oldUser)

	if oldUser.ID == 0 {
		return errors.New("User does not exist")
	}

	updatedUser.ID = oldUser.ID

	if user.Username != "" {
		updatedUser.Username = user.Username
	}

	if user.Password != "" {
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("couldn't hash password")
		}

		updatedUser.PasswordHash = newPasswordHash
	}

	if user.Bio != "" {
		updatedUser.Bio = user.Bio
	}

	if user.RoleUUID != "" {
		var role dbmodels.Role
		tools.DB.Table("roles").Select("roles.id").Where("roles.uuid = ?", user.RoleUUID).First(&role)
		if role.ID == 0 {
			return errors.New("Role does not exist")
		}

		updatedUser.RoleID = role.ID
	}

	tools.DB.Model(&updatedUser).Updates(updatedUser)
	return nil
}
