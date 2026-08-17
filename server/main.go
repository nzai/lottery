// 双色球走势图服务端：定时抓取福彩官网开奖数据，提供查询 API 并托管前端静态文件。
//
// 用法：
//
//	lottery-server            # 启动服务（首次自动全量回填，此后每日定时增量同步）
//	lottery-server -sync      # 手动触发一次同步后退出
package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/nzai/lottery/server/api"
	"github.com/nzai/lottery/server/config"
	"github.com/nzai/lottery/server/fetcher"
	"github.com/nzai/lottery/server/store"
)

// 前端构建产物（web 构建输出到 static/，部署时为单个二进制）。
//
//go:embed all:static
var staticFiles embed.FS

// buildTime 编译时间，部署时通过 ldflags 注入：
//
//	go build -ldflags "-X github.com/nzai/lottery/server.buildTime=2026-08-17T16:00:00Z"
//
// 未注入时为 "dev"（本地 go run 场景）。
var buildTime = "dev"

func main() {
	syncFlag := flag.Bool("sync", false, "手动触发一次同步后退出")
	flag.Parse()

	cfg := config.Load()

	st, err := store.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	defer st.Close()

	f := fetcher.New(cfg.UserAgent, cfg.FetchDelay)

	if *syncFlag {
		syncOnce(st, f)
		return
	}

	if cfg.FetchEnable {
		go backfill(st, f)
		scheduleSync(st, f, cfg.SyncCron)
	}

	// 嵌入的前端产物：剥掉 "static" 前缀后作为文件系统传给 API 层
	staticSub, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("读取嵌入的前端产物失败: %v", err)
	}

	router := api.NewRouter(st, staticSub, buildTime)
	log.Printf("HTTP 服务监听 %s", cfg.Addr)
	if err := router.Run(cfg.Addr); err != nil {
		log.Fatalf("HTTP 服务退出: %v", err)
	}
}

// backfill 启动时执行：库为空则全量回填，否则增量补最新。
func backfill(st *store.Store, f *fetcher.Fetcher) {
	count, err := st.Count()
	if err != nil {
		log.Printf("回填: 查询数量失败 %v", err)
		return
	}
	if count > 0 {
		syncOnce(st, f)
		return
	}
	log.Println("回填: 库为空，开始全量抓取（每页间隔抓取，可能需要数分钟）")
	start := time.Now()
	err = f.FetchAll(func(draws []store.Draw) error {
		return st.UpsertMany(draws)
	})
	if err != nil {
		log.Printf("回填失败: %v（服务继续运行，可稍后用 -sync 手动重试）", err)
		return
	}
	n, _ := st.Count()
	log.Printf("回填完成: 共 %d 期，耗时 %s", n, time.Since(start).Round(time.Second))
}

// syncOnce 抓取最新一页并幂等入库（每日增量同步 / -sync 手动同步共用）。
func syncOnce(st *store.Store, f *fetcher.Fetcher) {
	latest, err := st.LatestIssue()
	if err != nil {
		log.Printf("同步: 查询最新期号失败 %v", err)
		return
	}
	draws, err := f.FetchLatest()
	if err != nil {
		log.Printf("同步失败: %v", err)
		return
	}
	added := 0
	for _, d := range draws {
		if d.Issue > latest { // 期号定长，字符串比较即期次比较
			if err := st.Upsert(d); err != nil {
				log.Printf("同步: 写入 %s 失败 %v", d.Issue, err)
				return
			}
			added++
		}
	}
	log.Printf("同步完成: 拉取 %d 期，新增 %d 期", len(draws), added)
}

// scheduleSync 按 cron 表达式每日定时同步（Asia/Shanghai 时区，默认每天 21:30）。
func scheduleSync(st *store.Store, f *fetcher.Fetcher, expr string) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	c := cron.New(cron.WithLocation(loc))
	if _, err := c.AddFunc(expr, func() { syncOnce(st, f) }); err != nil {
		log.Printf("定时任务配置失败: %v", err)
		return
	}
	c.Start()
	log.Printf("定时同步已启用: %s（Asia/Shanghai）", expr)
}
