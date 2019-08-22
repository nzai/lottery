package util

import (
	"log"
	"time"

	"github.com/nzai/netop"
)

//	下载网页
func DownloadHtml(url string) ([]byte, error) {
	buffer, err := netop.GetBytes(url, netop.Retry(3, time.Second*2))
	if err != nil {
		log.Printf("获取网页%s失败: %v", url, err.Error())
		return nil, err
	}

	return buffer, err
}
