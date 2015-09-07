package util

import (
	"io/ioutil"
	"log"
	"net/http"
)

//	下载网页
func DownloadHtml(url string) ([]byte, error) {

	//  抓取网页
	response, err := http.Get(url)
	if err != nil {
		log.Println("获取网页失败: ", err.Error())
		return nil, err
	}
	defer response.Body.Close()

	//  读取网页
	buffer, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("读取网页失败: ", err.Error())
		return nil, err
	}

	return buffer, err
}
