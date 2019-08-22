package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nzai/lottery/config"
	"github.com/nzai/lottery/logic/job"
)

func main() {
	//	当前目录
	rootDir := filepath.Dir(os.Args[0])

	//	设置配置文件
	err := config.SetRootDir(rootDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	//	启动定时任务
	go job.Start()

	//	http监听端口
	port := config.Int("http", "port", 9000)
	serverAddress := fmt.Sprintf(":%d", port)

	r := gin.New()
	r.Use(gin.Logger())

	//	注册路由
	RegisterRoute(r)

	r.Run(serverAddress) // listen and serve on 0.0.0.0:8080
}
