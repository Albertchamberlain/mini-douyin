package serializer

import (
	"ADDD_DOUYIN/model"
	"ADDD_DOUYIN/util"
)

type UserLoginResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserRegisterResponse struct {
	Response
	UserId uint   `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserInfoResponse struct {
	Response
	*User
}

type UserList struct {
	Response
	UserList []*User `json:"user_list"`
}

func PackUser(u *model.User, userId uint, check, defaultTo bool) *User {
	if u == nil {
		return nil
	}

	followerCount, followeeCount, isFollow := int64(0), int64(0), defaultTo
	followerCount, _ = util.Redis.CountFollower(u.ID)
	followeeCount, _ = util.Redis.CountFollowee(u.ID)
	if check {
		isFollow, _ = util.Redis.IsFollow(userId, u.ID)
	}

	return &User{
		Id:            u.ID,
		Name:          u.Name,
		FollowCount:   followeeCount,
		FollowerCount: followerCount,
		IsFollow:      isFollow,
	}
}

func PackUsers(vs []*model.User, userId uint, check, defaultTo bool) []*User {
	users := make([]*User, 0)
	for _, v := range vs {
		if v2 := PackUser(v, userId, check, defaultTo); v2 != nil {
			users = append(users, v2)
		}
	}
	return users
}
