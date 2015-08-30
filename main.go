package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/nzai/lottery/config"
)

func main() {

	//	go func() {
	//		ticker := time.NewTicker(time.Second * 3)
	//		for _ = range ticker.C {
	//			log.Print("aa")
	//		}
	//	}()

	//	当前目录
	rootDir := filepath.Dir(os.Args[0])

	//	设置配置文件
	err := config.SetRootDir(rootDir)
	if err != nil {
		log.Fatal(err)
		return
	}

	//	http监听端口
	port := config.Int("http", "port", 9000)
	serverAddress := fmt.Sprintf(":%d", port)

	r := gin.New()
	r.Use(gin.Logger())

	//	静态文件目录
	r.Static("static", "./static")

	//	模板目录
	r.LoadHTMLGlob("static/html/*.html")

	r.GET("/", func(c *gin.Context) {
		c.String(200, "pong")
	})

	//	默认图标
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Redirect(http.StatusOK, "./static/icon/favicon.png")
	})
	
	//	双色球
	r.GET("/twocolorball", func(c *gin.Context) {
		
		c.HTML(http.StatusOK, "twocolorball/index.html", nil)
	})

	r.Run(serverAddress) // listen and serve on 0.0.0.0:8080
}
