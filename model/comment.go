package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	User    User
	UserId  uint
	VideoId uint
	Content string `gorm:"type:longtext"`
}
