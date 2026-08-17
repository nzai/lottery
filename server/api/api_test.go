package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"lottery/server/store"
)

func seedStore(t *testing.T) *store.Store {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	draws := []store.Draw{
		{Issue: "2026093", Date: "2026-08-13", Red: [6]int{5, 8, 15, 20, 21, 24}, Blue: 9},
		{Issue: "2026094", Date: "2026-08-16", Red: [6]int{6, 13, 15, 17, 24, 25}, Blue: 1},
	}
	if err := s.UpsertMany(draws); err != nil {
		t.Fatalf("UpsertMany: %v", err)
	}
	return s
}

func staticDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte(`<!doctype html><title>lottery</title>`), 0o644); err != nil {
		t.Fatalf("写 index.html: %v", err)
	}
	return dir
}

func TestHealth(t *testing.T) {
	router := NewRouter(seedStore(t), staticDir(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", w.Code)
	}
	if !json.Valid(w.Body.Bytes()) {
		t.Fatalf("响应不是合法 JSON: %s", w.Body.String())
	}
}

func TestDraws(t *testing.T) {
	router := NewRouter(seedStore(t), staticDir(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/draws", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", w.Code)
	}
	var body struct {
		Draws []struct {
			Issue string `json:"issue"`
			Date  string `json:"date"`
			Red   []int  `json:"red"`
			Blue  int    `json:"blue"`
		} `json:"draws"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("解析响应: %v", err)
	}
	if len(body.Draws) != 2 {
		t.Fatalf("draws 长度 = %d, want 2", len(body.Draws))
	}
	if body.Draws[0].Issue != "2026094" || body.Draws[0].Blue != 1 {
		t.Errorf("draws[0] = %+v, want 最新在前 2026094", body.Draws[0])
	}
	if len(body.Draws[0].Red) != 6 || body.Draws[0].Red[0] != 6 {
		t.Errorf("draws[0].Red = %v, want [6 13 15 17 24 25]", body.Draws[0].Red)
	}
}

func TestDrawsLimitAndBefore(t *testing.T) {
	router := NewRouter(seedStore(t), staticDir(t))

	// limit=1 → 只返回最新 1 期
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/draws?limit=1", nil)
	router.ServeHTTP(w, req)
	var body struct {
		Draws []struct {
			Issue string `json:"issue"`
		} `json:"draws"`
	}
	json.Unmarshal(w.Body.Bytes(), &body)
	if len(body.Draws) != 1 || body.Draws[0].Issue != "2026094" {
		t.Errorf("limit=1 结果 = %+v, want 仅 2026094", body.Draws)
	}

	// before=2026094 → 只返回更早的期
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/draws?before=2026094", nil)
	router.ServeHTTP(w, req)
	json.Unmarshal(w.Body.Bytes(), &body)
	if len(body.Draws) != 1 || body.Draws[0].Issue != "2026093" {
		t.Errorf("before 结果 = %+v, want 仅 2026093", body.Draws)
	}

	// 非法 limit → 回落默认 100
	w = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/draws?limit=abc", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("非法 limit code = %d, want 200", w.Code)
	}
}

func TestSPAFallback(t *testing.T) {
	router := NewRouter(seedStore(t), staticDir(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/some/spa/route", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("code = %d, want 200", w.Code)
	}
	if !contains(w.Body.String(), "lottery") {
		t.Errorf("SPA 回退未返回 index.html: %s", w.Body.String())
	}
}

func TestAPINotFound(t *testing.T) {
	router := NewRouter(seedStore(t), staticDir(t))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/nope", nil)
	router.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("code = %d, want 404", w.Code)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
