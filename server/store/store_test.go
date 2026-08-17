package store

import (
	"path/filepath"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func TestEmptyStore(t *testing.T) {
	s := newTestStore(t)
	latest, err := s.LatestIssue()
	if err != nil {
		t.Fatalf("LatestIssue: %v", err)
	}
	if latest != "" {
		t.Errorf("LatestIssue = %q, want 空串", latest)
	}
	n, err := s.Count()
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if n != 0 {
		t.Errorf("Count = %d, want 0", n)
	}
}

func TestUpsertListLatestCount(t *testing.T) {
	s := newTestStore(t)
	d1 := Draw{Issue: "2026093", Date: "2026-08-13", Red: [6]int{5, 8, 15, 20, 21, 24}, Blue: 9}
	d2 := Draw{Issue: "2026094", Date: "2026-08-16", Red: [6]int{6, 13, 15, 17, 24, 25}, Blue: 1}
	if err := s.Upsert(d1); err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if err := s.Upsert(d2); err != nil {
		t.Fatalf("Upsert: %v", err)
	}

	latest, _ := s.LatestIssue()
	if latest != "2026094" {
		t.Errorf("LatestIssue = %q, want 2026094", latest)
	}
	n, _ := s.Count()
	if n != 2 {
		t.Errorf("Count = %d, want 2", n)
	}

	// 幂等覆盖：重写 d2 的蓝球
	d2.Blue = 2
	if err := s.Upsert(d2); err != nil {
		t.Fatalf("Upsert 覆盖: %v", err)
	}
	n, _ = s.Count()
	if n != 2 {
		t.Errorf("Count = %d, want 2（覆盖不新增）", n)
	}

	draws, err := s.List(10, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(draws) != 2 {
		t.Fatalf("List 长度 = %d, want 2", len(draws))
	}
	if draws[0].Issue != "2026094" || draws[0].Blue != 2 {
		t.Errorf("draws[0] = %+v, want 2026094/blue=2（最新在前）", draws[0])
	}
	if draws[1].Issue != "2026093" {
		t.Errorf("draws[1].Issue = %q, want 2026093", draws[1].Issue)
	}

	// before 分页：只取更早的期
	older, err := s.List(10, "2026094")
	if err != nil {
		t.Fatalf("List before: %v", err)
	}
	if len(older) != 1 || older[0].Issue != "2026093" {
		t.Errorf("older = %+v, want 仅 2026093", older)
	}
}

func TestUpsertMany(t *testing.T) {
	s := newTestStore(t)
	draws := []Draw{
		{Issue: "2026090", Date: "2026-08-06", Red: [6]int{1, 2, 3, 4, 5, 6}, Blue: 1},
		{Issue: "2026091", Date: "2026-08-09", Red: [6]int{7, 8, 9, 10, 11, 12}, Blue: 2},
	}
	if err := s.UpsertMany(draws); err != nil {
		t.Fatalf("UpsertMany: %v", err)
	}
	n, _ := s.Count()
	if n != 2 {
		t.Errorf("Count = %d, want 2", n)
	}
	// 重复批量写入应幂等
	if err := s.UpsertMany(draws); err != nil {
		t.Fatalf("UpsertMany 重复: %v", err)
	}
	n, _ = s.Count()
	if n != 2 {
		t.Errorf("Count = %d, want 2（幂等）", n)
	}
}
