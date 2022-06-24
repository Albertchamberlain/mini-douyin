package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Author        User
	AuthorId      uint
	PlayUrl       string
	CoverUrl      string
	FavoriteCount int64
	CommentCount  int64
	Title         string
}
