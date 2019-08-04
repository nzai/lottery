package job

import (
	"log"
	"sync"
	"time"

	"github.com/nzai/lottery/config"
	"github.com/nzai/lottery/logic/superlotto"
	"github.com/nzai/lottery/logic/twocolorball"
)

func Start() {

	//	同步时间间隔
	intervalSeconds := config.Int("job", "intervalSeconds", 600)
	go func() {
		//	启动定时器
		ticker := time.NewTicker(time.Second * time.Duration(intervalSeconds))
		for _ = range ticker.C {
			syncData()
		}
	}()

	// 立即执行一次
	syncData()
}

func syncData() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		//  同步双色球
		log.Print("同步双色球开奖结果开始")
		err := twocolorball.SyncData()
		if err != nil {
			log.Print("同步双色球开奖结果失败: ", err)
		}
		log.Print("同步双色球开奖结果结束")

		wg.Done()
	}()

	go func() {
		//  同步大乐透
		log.Print("同步大乐透开奖结果开始")
		err := superlotto.SyncData()
		if err != nil {
			log.Print("同步大乐透开奖结果失败: ", err)
		}
		log.Print("同步大乐透开奖结果结束")

		wg.Done()
	}()

	wg.Wait()
}
