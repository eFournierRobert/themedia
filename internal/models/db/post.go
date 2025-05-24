package dbmodels

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UUID   string  `gorm:"type:char(36);uniqueIndex"`
	Title  *string `gorm:"type:varchar(255)"`
	Body   string
	UserID uint
	PostID *uint
}
