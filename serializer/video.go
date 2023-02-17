package serializer

import (
	"ADDD_DOUYIN/model"
	"ADDD_DOUYIN/util"
	"fmt"
)

type FeedResponse struct {
	Response
	VideoList []*Video `json:"video_list,omitempty"`
	NextTime  int64    `json:"next_time,omitempty"`
}

func PackVideo(v *model.Video, userId uint, check, defaultTo bool) *Video {
	if v == nil {
		return nil
	}

	favoriteCount, commentCount, isFavorite := int64(0), int64(0), defaultTo
	favoriteCount, _ = util.Redis.CountLiked(v.ID)
	commentCountInterface, _ := util.Redis.CountComment(v.ID)
	switch commentCountInterface.(type) {
	case int64:
		commentCount, _ = commentCountInterface.(int64)
	default:
		fmt.Println("assert fail")
	}
	if check {
		isFavorite, _ = util.Redis.IsLike(v.ID, userId)
	}

	return &Video{
		Id:            v.ID,
		Author:        *PackUser(&v.Author, userId, true, false),
		PlayUrl:       v.PlayUrl,
		CoverUrl:      v.CoverUrl,
		FavoriteCount: favoriteCount,
		CommentCount:  commentCount,
		IsFavorite:    isFavorite,
		Title:         v.Title,
	}
}

func PackVideos(vs []*model.Video, userId uint, check, defaultTo bool) []*Video {
	videos := make([]*Video, 0)
	for _, v := range vs {
		if v2 := PackVideo(v, userId, check, defaultTo); v2 != nil {
			videos = append(videos, v2)
		}
	}
	return videos
}
