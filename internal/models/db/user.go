package dbmodels

import "gorm.io/gorm"

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
