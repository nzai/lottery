// Package config 从环境变量加载服务配置，未设置时使用默认值。
package config

import (
	"os"
	"strconv"
	"time"
)

// Config 服务配置。
type Config struct {
	Addr        string        // HTTP 监听地址
	DBPath      string        // SQLite 文件路径
	SyncCron    string        // 每日同步 cron 表达式（Asia/Shanghai 时区）
	UserAgent   string        // 抓取请求 User-Agent
	FetchDelay  time.Duration // 分页抓取间隔
	FetchEnable bool          // 是否启用抓取与定时同步
}

// Load 从环境变量加载配置。
func Load() *Config {
	return &Config{
		Addr:        env("LOTTERY_ADDR", ":23817"),
		DBPath:      env("LOTTERY_DB", "lottery.db"),
		SyncCron:    env("LOTTERY_SYNC_CRON", "30 21 * * *"),
		UserAgent:   env("LOTTERY_UA", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"),
		FetchDelay:  time.Duration(envInt("LOTTERY_FETCH_DELAY_MS", 1500)) * time.Millisecond,
		FetchEnable: envBool("LOTTERY_FETCH_ENABLE", true),
	}
}

func env(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return def
}

func envBool(key string, def bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return def
}
