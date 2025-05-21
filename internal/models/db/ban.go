package dbmodels

import (
	"time"

	"gorm.io/gorm"
)

type Ban struct {
	gorm.Model
	UserId      uint
	EndDatetime time.Time
}
