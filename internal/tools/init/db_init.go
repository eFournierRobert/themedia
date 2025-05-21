package init_tools

import (
	"time"

	dbmodels "github.com/eFournierRobert/themedia/internal/models/db"
	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// startupDbMigration is the function responsible for doing
// a migration on the database at startup.
// If it isn't able to get access to the database, it will
// panic and print the error.
func StartupDbMigration() {
	err := tools.GetDb()
	if err != nil {
		panic(err.Error())
	}

	tools.DB.AutoMigrate(&dbmodels.Post{})
	tools.DB.AutoMigrate(&dbmodels.Role{})
	tools.DB.AutoMigrate(&dbmodels.User{})
	tools.DB.AutoMigrate(&dbmodels.Ban{})

	CheckIfFirstStartup(tools.DB)
}

// checkIfFirstStartup is the function that checks if the database
// has the required roles and user for deleted content. If not,
// it creates them.
func CheckIfFirstStartup(DB *gorm.DB) {
	var count int64
	DB.Model(&dbmodels.Role{}).Count(&count)

	if count == 0 {
		roles := []*dbmodels.Role{
			{Name: "admin", UUID: uuid.NewString()},
			{Name: "user", UUID: uuid.NewString()},
		}

		DB.Create(roles)
	}

	var user dbmodels.User
	DB.Select("id").Where("username = ?", "deleted").First(&user)

	if user.ID == 0 {
		DB.Create(&dbmodels.User{
			Username: "deleted",
			UUID:     uuid.NewString(),
			RoleID:   2,
		})

		DB.Create(&dbmodels.Ban{
			UserId:      1,
			EndDatetime: time.Now().AddDate(1000, 0, 0),
		})
	}
}
