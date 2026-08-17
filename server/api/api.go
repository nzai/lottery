// Package api 提供 HTTP 接口与前端静态文件托管。
package api

import (
	"io/fs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/nzai/lottery/server/store"
)

// NewRouter 构建路由：数据 API + 静态托管（嵌入 fs）+ SPA 回退。
// staticFS 为前端构建产物的文件系统（main 包通过 go:embed 传入）。
func NewRouter(s *store.Store, staticFS fs.FS) *gin.Engine {
	r := gin.Default()

	// 静态资源缓存策略：
	// - /assets/* 文件名带内容 hash，内容不变 hash 不变 → 一年 immutable，浏览器不再发请求
	// - /（index.html）不带 hash → no-cache 每次校验，bug 修复后用户立即拿到新入口
	r.Use(func(c *gin.Context) {
		p := c.Request.URL.Path
		switch {
		case strings.HasPrefix(p, "/assets/"):
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		case p == "/" || p == "/index.html":
			c.Header("Cache-Control", "no-cache")
		}
		c.Next()
	})

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

	// 前端静态文件（嵌入的构建产物）。
	// 注意：
	// - 不能用 http.FileServer 处理 index.html：标准库会把以 "/index.html" 结尾的
	//   路径 301 重定向到 "./"，所以 index.html/favicon 直接用 fs.ReadFile 返回
	// - assets 挂载到 fs 的 assets 子目录（URL /assets/xxx ↔ fs assets/xxx）
	assetsFS := staticFS
	if sub, err := fs.Sub(staticFS, "assets"); err == nil {
		assetsFS = sub // 未构建时无 assets 目录，fallback 后自然 404
	}
	assetsServer := http.StripPrefix("/assets/", http.FileServer(http.FS(assetsFS)))
	r.GET("/assets/*filepath", func(c *gin.Context) {
		assetsServer.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/", func(c *gin.Context) {
		serveEmbedded(c, staticFS, "index.html", "text/html; charset=utf-8")
	})
	r.GET("/favicon.svg", func(c *gin.Context) {
		serveEmbedded(c, staticFS, "favicon.svg", "image/svg+xml")
	})

	// SPA 回退：非 /api 的未知路径返回 index.html
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		serveEmbedded(c, staticFS, "index.html", "text/html; charset=utf-8")
	})
	return r
}

// serveEmbedded 从嵌入 fs 读取文件并返回；不存在则 404。
func serveEmbedded(c *gin.Context, fsys fs.FS, name, contentType string) {
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	c.Data(http.StatusOK, contentType, data)
}
