package tools

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	UUID string `gorm:"type:char(36);uniqueIndex"`
	Name string
	Users []User
}

type User struct {
	gorm.Model
	UUID string `gorm:"type:char(36);uniqueIndex"`
	Username string
	PasswordHash []byte
	RoleID uint
}

type FullUser struct {
	ID uint
	UUID string
	Username string
	Role Role
}

func CreateUser(username *string, password *string, role *Role) (*User, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Couldn't hash password")
	}

	if role.ID == 0 {
		db.Where("name = ?", "user").First(&role)
		if role.ID == 0 {
			return nil, errors.New("Couldn't find the user role")
		}
	}

	var user User

	user = User {
		UUID: uuid.NewString(),
		Username: *username,
		PasswordHash: hashedPassword,
		RoleID: role.ID,
	}
		

	db.Create(&user)

	return &user, nil
}

func FindFullUserByUUID(uuid *string) (*User, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	var user User
	db.Table("users").Select(
		"user.id", 
		"user.uuid", 
		"user.username", 
		"role.uuid", 
		"role.name",
	).Where("uuid = ?", *uuid).Group("roles").First(&user)

	return &user, nil
}

func FindRoleByID(id int) (*Role, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	var role Role
	db.Where("id = ?", id).First(&role)


	return &role, nil
}

func FindRoleByUUID(uuid *string) (*Role, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	var role Role
	db.Where("uuid = ?", *uuid).First(&role)


	return &role, nil
}
