package main

import (
	"ADDD_DOUYIN/conf"
	"ADDD_DOUYIN/routes"
	"fmt"
	_ "net/http/pprof"
)

func main() {
	conf.Init()
	r := routes.NewRouter()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("服务器启动发生错误,信息如下:", err)
	}
}
