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

func FindRole(uuid *string) (*Role, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	var role Role
	db.Where("uuid = ?", *uuid).First(&role)


	return &role, nil
}
