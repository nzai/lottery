// Package api 提供 HTTP 接口与前端静态文件托管。
package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"lottery/server/store"
)

// NewRouter 构建路由：数据 API + 静态托管 + SPA 回退。
func NewRouter(s *store.Store, staticDir string) *gin.Engine {
	r := gin.Default()

	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/api/draws", func(c *gin.Context) {
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
		if err != nil || limit <= 0 || limit > 500 {
			limit = 100
		}
		draws, err := s.List(limit, c.Query("before"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"draws": draws})
	})

	// 前端静态文件（构建产物输出到 staticDir）。
	// 注意不能用 StaticFS("/")：它的 /*filepath 通配路由与 /api/* 冲突。
	r.Static("/assets", staticDir+"/assets")
	r.StaticFile("/", staticDir+"/index.html")
	r.StaticFile("/favicon.svg", staticDir+"/favicon.svg")

	// SPA 回退：非 /api 的未知路径返回 index.html
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.File(staticDir + "/index.html")
	})
	return r
}
