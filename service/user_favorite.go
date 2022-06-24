package service

import (
	"ADDD_DOUYIN/conf"
	"ADDD_DOUYIN/model"
	"ADDD_DOUYIN/util"
)

type UserFavoriteAction struct {
	UserId     uint   `form:"user_id"`
	Token      string `form:"token"`
	VideoId    uint   `form:"video_id"`
	ActionType int    `form:"action_type"`
}

func (f *UserFavoriteAction) Action() error {
	switch f.ActionType {
	case 1:
		return f.create()
	case 2:
		return f.delete()
	}
	return nil
}

func (f *UserFavoriteAction) create() error {
	return util.Redis.Like(f.UserId, f.VideoId)
}

func (f *UserFavoriteAction) delete() error {
	return util.Redis.Unlike(f.UserId, f.VideoId)
}

func FavoriteList(userId uint) ([]*model.Video, error) {
	videoIdList, err := util.Redis.RangeLike(userId, 0, -1)
	if err != nil {
		return nil, err
	}
	videos := make([]*model.Video, len(videoIdList))
	if len(videoIdList) == 0 {
		return videos, err
	}
	if err = conf.DB.Model(&model.Video{}).Find(&videos, videoIdList).Error; err != nil {
		return nil, err
	}
	return videos, err
}
