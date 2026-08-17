package fetcher

import (
	"testing"
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
