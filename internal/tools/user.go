package tools

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UUID string
	Username string
	PasswordHash []byte
}

func CreateUser(username string, password string) (*User, error) {
	db, err := GetDb()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Couldn't hash password")
	}

	var user User = User {
		UUID: uuid.NewString(),
		Username: username,
		PasswordHash: hashedPassword,
	}

	db.Create(&user)

	return &user, nil
}
