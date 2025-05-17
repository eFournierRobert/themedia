package ban_tools

import (
	"errors"
	"time"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	"github.com/eFournierRobert/themedia/internal/tools"
)

// CreateBan is the function that inserts a ban in the database.
// It will return an error if the user does not exist.
func CreateBan(userUUID string, endDatetime time.Time) error {
	var user dbmodels.User
	tools.DB.Table("users").Select("id").Where("uuid = ?", userUUID).First(&user)
	if user.ID == 0 {
		return errors.New("user does not exist")
	}

	tools.DB.Create(&dbmodels.Ban{
		UserId:      user.ID,
		EndDatetime: endDatetime,
	})

	return nil
}

// IsUserBanned is the function that checks if the given user
// is currently banned. Will return true if yes and false if not.
func IsUserBanned(userUUID string) bool {
	var ban dbmodels.Ban
	tools.DB.Where("end_datetime >= ?", time.Now()).First(&ban)

	return ban.ID != 0
}
