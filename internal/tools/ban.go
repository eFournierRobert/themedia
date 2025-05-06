package tools

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Ban struct {
	gorm.Model
	UserId      uint
	EndDatetime time.Time
}

// CreateBan is the function that inserts a ban in the database.
// It will return an error if the user does not exist.
func CreateBan(userUUID string, endDatetime time.Time) error {
	var user User
	DB.Table("users").Select("id").Where("uuid = ?", userUUID).First(&user)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}

	DB.Create(&Ban{
		UserId:      user.ID,
		EndDatetime: endDatetime,
	})

	return nil
}

// IsUserBanned is the function that checks if the given user
// is currently banned. Will return true if yes and false if not.
func IsUserBanned(userUUID string) bool {
	var ban Ban
	DB.Where("end_datetime >= ?", time.Now()).First(&ban)

	return ban.ID != 0
}
