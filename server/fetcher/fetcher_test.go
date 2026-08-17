package fetcher

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"lottery/server/store"
)

func TestParseDraw(t *testing.T) {
	d, err := ParseDraw("2026094", "2026-08-16(日)", "06,13,15,17,24,25", "01")
	if err != nil {
		t.Fatalf("ParseDraw: %v", err)
	}
	if d.Issue != "2026094" {
		t.Errorf("Issue = %q, want 2026094", d.Issue)
	}
	if d.Date != "2026-08-16" {
		t.Errorf("Date = %q, want 2026-08-16（去掉星期后缀）", d.Date)
	}
	wantRed := [6]int{6, 13, 15, 17, 24, 25}
	if d.Red != wantRed {
		t.Errorf("Red = %v, want %v（零填充转数字）", d.Red, wantRed)
	}
	if d.Blue != 1 {
		t.Errorf("Blue = %d, want 1", d.Blue)
	}
}

func TestParseDrawErrors(t *testing.T) {
	cases := []struct {
		name                 string
		code, date, red, blue string
	}{
		{"红球个数不足", "2026094", "2026-08-16", "06,13,15,17,24", "01"},
		{"红球个数超限", "2026094", "2026-08-16", "06,13,15,17,24,25,26", "01"},
		{"红球越界", "2026094", "2026-08-16", "06,13,15,17,24,34", "01"},
		{"蓝球越界", "2026094", "2026-08-16", "06,13,15,17,24,25", "17"},
		{"蓝球非法字符", "2026094", "2026-08-16", "06,13,15,17,24,25", "ab"},
		{"日期非法", "2026094", "2026-16-16", "06,13,15,17,24,25", "01"},
		{"空日期", "2026094", "", "06,13,15,17,24,25", "01"},
	}
	for _, c := range cases {
		if _, err := ParseDraw(c.code, c.date, c.red, c.blue); err == nil {
			t.Errorf("%s: 期望报错，实际成功", c.name)
		}
	}
}

// loadFixture 读取官方真实返回样例。
func loadFixture(t *testing.T) string {
	t.Helper()
	b, err := os.ReadFile("testdata/page1.json")
	if err != nil {
		t.Fatalf("读取 fixture: %v", err)
	}
	return string(b)
}

func TestFetchPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "ssq" {
			t.Errorf("name 参数 = %q, want ssq", r.URL.Query().Get("name"))
		}
		if r.Header.Get("User-Agent") == "" {
			t.Error("缺少 User-Agent")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(loadFixture(t)))
	}))
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	f := New("test-ua", 0)
	draws, err := f.FetchPage(1, 100)
	if err != nil {
		t.Fatalf("FetchPage: %v", err)
	}
	if len(draws) != 2 {
		t.Fatalf("len = %d, want 2", len(draws))
	}
	if draws[0].Issue != "2026094" || draws[0].Blue != 1 {
		t.Errorf("draws[0] = %+v, want 2026094/blue=1", draws[0])
	}
	if draws[1].Issue != "2026093" {
		t.Errorf("draws[1].Issue = %q, want 2026093", draws[1].Issue)
	}
}

func TestFetchPageError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	f := New("test-ua", 0)
	if _, err := f.FetchPage(1, 100); err == nil {
		t.Error("HTTP 403 应报错")
	}
}

func TestFetchAllPages(t *testing.T) {
	// total=2052, pageSize=100 → 21 页；模拟按页返回，校验页数与回调顺序
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pageNo := r.URL.Query().Get("pageNo")
		body := loadFixture(t)
		// 第 2 页起替换期号为 2026092、2026091，便于校验
		if pageNo == "2" {
			body = `{"state":0,"message":"ok","total":2052,"result":[{"code":"2026092","date":"2026-08-09(日)","red":"01,02,03,04,05,06","blue":"02"},{"code":"2026091","date":"2026-08-06(四)","red":"11,12,13,14,15,16","blue":"03"}]}`
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
		page++
	}))
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	f := New("test-ua", 0)
	var pages, draws int
	err := f.FetchAll(func(ds []store.Draw) error {
		pages++
		draws += len(ds)
		return nil
	})
	if err != nil {
		t.Fatalf("FetchAll: %v", err)
	}
	if pages != 21 {
		t.Errorf("pages = %d, want 21（total 2052 / 100 向上取整）", pages)
	}
	if draws != 21*2 {
		t.Errorf("draws = %d, want %d", draws, 21*2)
	}
}

func TestFetchLatest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("pageNo") != "1" {
			t.Errorf("pageNo = %q, want 1", r.URL.Query().Get("pageNo"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(loadFixture(t)))
	}))
	defer srv.Close()

	oldBase := baseURL
	baseURL = srv.URL
	defer func() { baseURL = oldBase }()

	f := New("test-ua", 0)
	draws, err := f.FetchLatest()
	if err != nil {
		t.Fatalf("FetchLatest: %v", err)
	}
	if len(draws) != 2 {
		t.Fatalf("len = %d, want 2", len(draws))
	}
}
