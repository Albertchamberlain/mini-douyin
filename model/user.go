package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name          string `gorm:"unique"`
	Password      string
	FollowCount   int64
	FollowerCount int64
}

func (user *User) SetPassword(password string) {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //非明文存储密码
	user.Password = string(bytes)
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
