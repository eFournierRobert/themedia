// Package tools is the package containing
// all the request made to the database.
package tools

import (
	"errors"

	"github.com/eFournierRobert/themedia/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role is the struct responsible for the table roles in the database.
type Role struct {
	gorm.Model
	UUID  string `gorm:"type:char(36);uniqueIndex"`
	Name  string
	Users []User
}

// User is the struct responsible for the table users in the database.
type User struct {
	gorm.Model
	UUID         string `gorm:"type:char(36);uniqueIndex"`
	Username     string
	PasswordHash []byte
	RoleID       uint
	Bio          string
	Bans         []Ban
}

// FullUser is the struct responsible to store the return value
// of a SELECT in the database that has the user information
// and the role information of that users.
type FullUser struct {
	ID       uint
	UserUUID string
	Username string
	Bio      string
	RoleUUID string
	Name     string
}

// CreateUser is the function responsible for inserting a new user in the database.
// the password passed must be in plain text since function will be hashing.
// It returns a pointer to a User struct or an error.
func CreateUser(username *string, password *string, role *Role) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Couldn't hash password")
	}

	if role.ID == 0 {
		DB.Where("name = ?", "user").First(&role)
		if role.ID == 0 {
			return nil, errors.New("Couldn't find the user role")
		}
	}

	var user User

	user = User{
		UUID:         uuid.NewString(),
		Username:     *username,
		PasswordHash: hashedPassword,
		RoleID:       role.ID,
	}

	DB.Create(&user)

	return &user, nil
}

// FindFullUserByUUID is the function responsible for finding the user that
// has the given UUID in the database.
// It will return a pointer to a FullUser struct or an error.
func FindFullUserByUUID(uuid *string) (*FullUser, error) {
	var fullUser FullUser
	DB.Table("users").Select(
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
func FindRoleByUUID(uuid *string) (*Role, error) {
	var role Role
	DB.Where("uuid = ?", *uuid).First(&role)

	return &role, nil
}

// VerifyPassword is the function that receives a raw passwords and
// compares it to the hashed one in the database for the user with the given
// UUID. It will return a boolean or an error.
func VerifyPassword(uuid *string, password *string) (bool, error) {
	var user User
	DB.Table("users").Select("id", "password_hash").Where("uuid = ?", uuid).First(&user)

	if user.ID == 0 {
		return false, errors.New("Did not find user")
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
	var user User
	DB.Table("users").Select("id").Where("uuid = ?", uuid).First(&user)

	return user.ID != 0
}

// IsUserAdmin takes the given user UUID and checks if the user
// is an admin or not.
func IsUserAdmin(uuid string) bool {
	var user FullUser
	DB.Table("users").Select(
		"users.id",
		"roles.name",
	).Where("users.uuid = ?", uuid).Joins("JOIN roles on roles.id = users.role_id").First(&user)

	return user.ID != 0 && user.Name == "admin"
}

// DeleteUser takes the given user UUID and hard delete it
// from the database.
func DeleteUser(uuid string) error {
	var user User
	DB.Where("uuid = ?", uuid).First(&user)

	if user.ID == 0 {
		return errors.New("User does not exist")
	}

	DB.Unscoped().Delete(&user)
	return nil
}

// UpdateUser updates the user in the database with the new information received.
// Returns nil if successful or an error if not.
func UpdateUser(uuid string, user *models.UserPost) error {
	var oldUser User
	var updatedUser User
	DB.Table("users").Select("id").Where("uuid = ?", uuid).First(&oldUser)

	updatedUser.ID = oldUser.ID

	if user.Username != "" {
		updatedUser.Username = user.Username
	}

	if user.Password != "" {
		newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(*&user.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("Couldn't hash password")
		}

		updatedUser.PasswordHash = newPasswordHash
	}

	if user.Bio != "" {
		updatedUser.Bio = user.Bio
	}

	if user.RoleUUID != "" {
		var role Role
		DB.Table("roles").Select("roles.id").Where("roles.uuid = ?", user.RoleUUID).First(&role)
		if role.ID == 0 {
			return errors.New("Role does not exist")
		}

		updatedUser.RoleID = role.ID
	}

	DB.Model(&updatedUser).Updates(updatedUser)
	return nil
}
