// Package fetcher 抓取并解析福彩官网双色球开奖数据。
package fetcher

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"lottery/server/store"
)

// 官方查询接口。
const baseURL = "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice"

// apiResponse 官方接口返回结构（仅取用到的字段，其余忽略）。
type apiResponse struct {
	State   int    `json:"state"` // 0 表示成功
	Message string `json:"message"`
	Total   int    `json:"total"` // 结果总数
	Result  []struct {
		Code string `json:"code"` // 期号
		Date string `json:"date"` // "2026-08-16(日)"
		Red  string `json:"red"`  // "06,13,15,17,24,25"
		Blue string `json:"blue"` // "01"
	} `json:"result"`
}

// ParseDraw 解析官方接口单条记录。
// date 形如 "2026-08-16(日)"，red 形如 "06,13,15,17,24,25"，blue 形如 "01"。
func ParseDraw(code, date, redStr, blueStr string) (store.Draw, error) {
	if i := strings.IndexByte(date, '('); i >= 0 {
		date = date[:i]
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return store.Draw{}, fmt.Errorf("日期 %q 非法", date)
	}
	red, err := parseNumbers(redStr, 6, 1, 33)
	if err != nil {
		return store.Draw{}, err
	}
	blue, err := parseNumbers(blueStr, 1, 1, 16)
	if err != nil {
		return store.Draw{}, err
	}
	return store.Draw{
		Issue: code,
		Date:  date,
		Red:   [6]int{red[0], red[1], red[2], red[3], red[4], red[5]},
		Blue:  blue[0],
	}, nil
}

// parseNumbers 解析逗号分隔的号码串，校验个数与范围。
func parseNumbers(s string, want, min, max int) ([]int, error) {
	parts := strings.Split(s, ",")
	if len(parts) != want {
		return nil, fmt.Errorf("号码个数 %d != %d", len(parts), want)
	}
	nums := make([]int, 0, want)
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || n < min || n > max {
			return nil, fmt.Errorf("号码 %q 非法", p)
		}
		nums = append(nums, n)
	}
	return nums, nil
}
