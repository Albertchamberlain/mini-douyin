package service

import (
	"ADDD_DOUYIN/conf"
	"ADDD_DOUYIN/model"
	"ADDD_DOUYIN/serializer"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type FeedService struct { //TODO 目前用户登录与否都可以返回Feed流，或许以后可以根据用户的特征返回一些定制化的Feed流
	LatestTime string `json:"latest_time,omitempty"`
	Token      string `json:"token,omitempty"`
}

func Feed(latestTime string) ([]*model.Video, error) {
	vs := make([]*model.Video, 30)
	// fixme filter by latestTime
	err := conf.DB.Preload("Author").Find(&vs).Order("created_at DESC").Limit(30).Error
	return vs, err
}

func (service *FeedService) FeedService() serializer.FeedResponse {
	videoList := make([]*serializer.Video, 30)
	if err := conf.DB.Select("*").Order("created_at DESC").Limit(30).Scan(&videoList).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(err)
			return serializer.FeedResponse{
				Response: serializer.Response{StatusCode: 1,
					StatusMsg: "内容不足30条,请稍后再来",
				},
			}
		}
		fmt.Println(err)
		return serializer.FeedResponse{
			Response: serializer.Response{StatusCode: 1,
				StatusMsg: "数据库出错",
			},
		}
	}
	return serializer.FeedResponse{
		Response: serializer.Response{StatusCode: 0,
			StatusMsg: "feed流返回成功",
		},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	}
}
