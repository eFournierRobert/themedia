package tools

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// startupDbMigration is the function responsible for doing
// a migration on the database at startup.
// If it isn't able to get access to the database, it will
// panic and print the error.
func StartupDbMigration() {
	err := GetDb()
	if err != nil {
		panic(err.Error())
	}

	DB.AutoMigrate(&Role{})
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Ban{})

	checkIfFirstStartup(DB)
}

// checkIfFirstStartup is the function that checks if the database
// has the required roles and user for deleted content. If not,
// it creates them.
func checkIfFirstStartup(DB *gorm.DB) {
	var count int64
	DB.Model(&Role{}).Count(&count)

	if count == 0 {
		roles := []*Role{
			{Name: "admin", UUID: uuid.NewString()},
			{Name: "user", UUID: uuid.NewString()},
		}

		DB.Create(roles)
	}

	var user User
	DB.Select("id").Where("username = ?", "deleted").First(&user)

	if user.ID == 0 {
		DB.Create(&User{
			Username: "deleted",
			UUID:     uuid.NewString(),
			RoleID:   2,
		})

		DB.Create(&Ban{
			UserId:      1,
			EndDatetime: time.Now().AddDate(1000, 0, 0),
		})
	}
}
