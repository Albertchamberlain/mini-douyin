package controller

import (
	"ADDD_DOUYIN/serializer"
	"ADDD_DOUYIN/service"
	"ADDD_DOUYIN/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponseForSwagger struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
	NextTime   int    `json:"next_time"`
	VideoList  []struct {
		ID     int `json:"id"`
		Author struct {
			ID            int    `json:"id"`
			Name          string `json:"name"`
			FollowCount   int    `json:"follow_count"`
			FollowerCount int    `json:"follower_count"`
			IsFollow      bool   `json:"is_follow"`
		} `json:"author"`
		PlayURL       string `json:"play_url"`
		CoverURL      string `json:"cover_url"`
		FavoriteCount int    `json:"favorite_count"`
		CommentCount  int    `json:"comment_count"`
		IsFavorite    bool   `json:"is_favorite"`
		Title         string `json:"title"`
	} `json:"video_list"`
}

// @Tags 用户相关接口
// @Summary 用户推送Feeds
// @Description 给游客或者注册用户推送Feeds的接口，一次30条
// @Router /feed [get]
// @Param latest_time query string false "上次访问时间"
// @Param token query string  false "token"
// @Produce json
// @Success 200 {object} FeedResponseForSwagger
func Feed(c *gin.Context) {

	latestTime := c.Query("latest_time")
	if latestTime == "" {
		latestTime = strconv.FormatInt(time.Now().Unix(), 10)
	}

	// 如果登陆了 对于feed流检查like情况 不登陆时默认为false
	userId, check := uint(0), false
	token, claim, err := util.ParseToken(c.Query("token"))
	if err == nil && token.Valid {
		userId = claim.Id
		check = true
	}

	if res, err := service.Feed(latestTime); err != nil {
		c.JSON(http.StatusOK, serializer.ConvertErr(err))
	} else {
		c.JSON(http.StatusOK, serializer.FeedResponse{
			Response:  serializer.Success,
			VideoList: serializer.PackVideos(res, userId, check, false),
			NextTime:  time.Now().Unix() / 1000,
		})
	}
}
