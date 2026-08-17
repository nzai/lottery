// Package fetcher 抓取并解析福彩官网双色球开奖数据。
package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/nzai/lottery/server/store"
)

// 官方查询接口（var 便于测试替换为 mock 服务器）。
var baseURL = "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice"

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

// Fetcher 负责分页抓取，带 User-Agent 与抓取间隔。
type Fetcher struct {
	client *http.Client
	ua     string
	delay  time.Duration
}

// New 创建抓取器。delay 为分页抓取间隔，用于礼貌抓取（测试传 0）。
func New(ua string, delay time.Duration) *Fetcher {
	return &Fetcher{
		client: &http.Client{Timeout: 30 * time.Second},
		ua:     ua,
		delay:  delay,
	}
}

// FetchPage 抓取第 pageNo 页（从 1 开始），返回该页数据（官方按开奖日期倒序）。
func (f *Fetcher) FetchPage(pageNo, pageSize int) ([]store.Draw, error) {
	raw, err := f.fetchRaw(pageNo, pageSize)
	if err != nil {
		return nil, err
	}
	return pageToDraws(raw)
}

// FetchLatest 抓取最新一页（每日增量同步用）。
func (f *Fetcher) FetchLatest() ([]store.Draw, error) {
	return f.FetchPage(1, 100)
}

// FetchAll 分页抓取全量历史数据，每抓一页回调一次（最新页在前）。
func (f *Fetcher) FetchAll(onPage func([]store.Draw) error) error {
	const pageSize = 100
	first, err := f.fetchRaw(1, pageSize)
	if err != nil {
		return err
	}
	pages := (first.Total + pageSize - 1) / pageSize

	if err := onPage(mustDraws(first)); err != nil {
		return err
	}
	for page := 2; page <= pages; page++ {
		time.Sleep(f.delay)
		raw, err := f.fetchRaw(page, pageSize)
		if err != nil {
			return fmt.Errorf("抓取第 %d/%d 页失败: %w", page, pages, err)
		}
		if err := onPage(mustDraws(raw)); err != nil {
			return err
		}
	}
	return nil
}

// fetchRaw 请求一页并解析 JSON。
func (f *Fetcher) fetchRaw(pageNo, pageSize int) (*apiResponse, error) {
	u := fmt.Sprintf("%s?name=ssq&pageNo=%d&pageSize=%d&systemType=PC", baseURL, pageNo, pageSize)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", f.ua)
	req.Header.Set("Referer", "https://www.cwl.gov.cn/")
	resp, err := f.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("接口返回 HTTP %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var ar apiResponse
	if err := json.Unmarshal(body, &ar); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}
	if ar.State != 0 {
		return nil, fmt.Errorf("接口返回异常 state=%d message=%s", ar.State, ar.Message)
	}
	return &ar, nil
}

// pageToDraws 将一页原始数据转为 Draw 列表；单条非法即报错。
func pageToDraws(raw *apiResponse) ([]store.Draw, error) {
	draws := make([]store.Draw, 0, len(raw.Result))
	for _, r := range raw.Result {
		d, err := ParseDraw(r.Code, r.Date, r.Red, r.Blue)
		if err != nil {
			return nil, fmt.Errorf("解析期号 %s 失败: %w", r.Code, err)
		}
		draws = append(draws, d)
	}
	return draws, nil
}

// mustDraws 是 pageToDraws 的封装：FetchAll 回调场景错误即整体失败。
func mustDraws(raw *apiResponse) []store.Draw {
	draws, err := pageToDraws(raw)
	if err != nil {
		panic(err) // FetchAll 内部使用，调用方已统一返回错误
	}
	return draws
}
