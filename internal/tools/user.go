// Package tools is the package containing
// all the request made to the database.
package tools

import (
	"errors"

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
}

// FullUser is the struct responsible to store the return value
// of a SELECT in the database that has the user information
// and the role information of that users.
type FullUser struct {
	ID       uint
	UserUUID string
	Username string
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
