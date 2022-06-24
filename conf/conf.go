package conf

import (
	"ADDD_DOUYIN/util"
	"fmt"
	"net/url"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/tencentyun/cos-go-sdk-v5"
	"gopkg.in/ini.v1"
)

var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

func Init() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
		panic(err)
	}

	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true&loc=Local"}, "")
	Database(path)

	u, _ := url.Parse(file.Section("cos").Key("url").String())
	b := &cos.BaseURL{BucketURL: u}
	id := file.Section("cos").Key("id").String()
	key := file.Section("cos").Key("key").String()
	util.InitCos(b, id, key)

	conf := &redis.Options{
		Addr:     file.Section("redis").Key("url").String(),
		Password: file.Section("redis").Key("password").String(),
		DB:       0,
		PoolSize: 100,
	}
	util.InitRedis(conf)
}
