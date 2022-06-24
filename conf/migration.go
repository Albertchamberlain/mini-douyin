package conf

import (
	"ADDD_DOUYIN/model"
)

func migration() {
	if err := DB.AutoMigrate(
		&model.User{},
		&model.Video{},
		&model.Comment{},
		&model.UserFollow{},
		&model.UserFavorite{},
	); err != nil {
		panic(err)
	}
}
