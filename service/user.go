package service

import (
	"ADDD_DOUYIN/conf"
	"ADDD_DOUYIN/model"
	"ADDD_DOUYIN/serializer"
	"ADDD_DOUYIN/util"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

//接受前端传来的用户登录
type UserService struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

//接受前端传来的用户信息
type UserInfoService struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

//UserRegisterService 用户注册服务
func (service *UserService) Register(username, password string) error {
	var user model.User
	if (util.Filter.KeySize() != 0 && util.Filter.TestString(username)) || !errors.Is(conf.DB.Where("name = ?", username).First(&user).Error, gorm.ErrRecordNotFound) {
		//fmt.Println(gorm.ErrRecordNotFound)
		return gorm.ErrRecordNotFound
	}
	user.Name = username
	//user.ID = util.NextId()  //已使用数据库自增ID主键
	user.SetPassword(password)
	//创建用户
	if err := conf.DB.Create(&user).Error; err != nil {
		//fmt.Println(err)
		return err
	}
	token, err := util.GenerateToken(user.ID, service.UserName, 0)
	if err != nil {
		//fmt.Println(err)
		return err
	}
	util.Filter.AddString(username) //将已注册的username写入bloomfilter
	return &serializer.UserRegisterResponse{
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  "注册成功",
		},
		UserId: user.ID,
		Token:  token,
	}
}

//Login 用户登陆服务
func (service *UserService) Login(username, password string) serializer.UserLoginResponse {
	var user model.User
	if err := conf.DB.Where("name=?", username).First(&user).Error; err != nil {
		//如果查询不到，返回相应的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(err)
			return serializer.UserLoginResponse{
				Response: serializer.Response{StatusCode: 1,
					StatusMsg: "用户不存在",
				},
			}
		}
		fmt.Println(err)
		return serializer.UserLoginResponse{
			Response: serializer.Response{StatusCode: 1,
				StatusMsg: "数据库出错",
			},
		}
	}
	if !user.CheckPassword(password) {
		return serializer.UserLoginResponse{
			Response: serializer.Response{StatusCode: 1,
				StatusMsg: "密码不匹配",
			},
		}
	}
	token, err := util.GenerateToken(user.ID, user.Name, 0)
	if err != nil {
		fmt.Println(err)
		return serializer.UserLoginResponse{
			Response: serializer.Response{StatusCode: 1,
				StatusMsg: "生成token出错",
			},
		}
	}
	return serializer.UserLoginResponse{ //登录成功
		Response: serializer.Response{StatusCode: 0,
			StatusMsg: "登录成功",
		},
		UserId: user.ID,
		Token:  token,
	}
}

//用户信息服务
func (service *UserInfoService) UserInfo(userId uint) serializer.UserInfoResponse {
	var user model.User
	if err := conf.DB.Where("id=?", userId).First(&user).Error; err != nil {
		//如果查询不到，返回相应的错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println(err)
			return serializer.UserInfoResponse{
				Response: serializer.Response{StatusCode: 1,
					StatusMsg: "用户不存在",
				},
			}
		}
		fmt.Println(err)
		return serializer.UserInfoResponse{
			Response: serializer.Response{StatusCode: 1,
				StatusMsg: "数据库出错",
			},
		}
	}
	//来到正常的处理逻辑
	return serializer.UserInfoResponse{ //正常请求返回用户信息
		Response: serializer.Response{
			StatusCode: 0,
			StatusMsg:  "用户信息查询成功",
		},
		User: serializer.PackUser(&user, userId, false, false),
	}
}

func Publish(video *model.Video) error {
	return conf.DB.Create(&video).Error
}

func PublishList(userId uint) ([]*model.Video, error) {
	videos := make([]*model.Video, 0)
	err := conf.DB.Where("author_id = ?", userId).Find(&videos).Error
	return videos, err
}
