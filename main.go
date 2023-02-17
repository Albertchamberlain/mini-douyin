package main

import (
	"ADDD_DOUYIN/conf"
	_ "ADDD_DOUYIN/docs"
	"ADDD_DOUYIN/routes"
	"fmt"
	_ "net/http/pprof"

	"github.com/gin-gonic/gin"
)

func main() {
	conf.Init()
	gin.ForceConsoleColor()
	r := routes.NewRouter()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("服务器启动发生错误,信息如下:", err)
	}
}
