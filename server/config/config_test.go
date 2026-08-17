package config

import (
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	if cfg.Addr != ":23817" {
		t.Errorf("Addr = %q, want :23817", cfg.Addr)
	}
	if cfg.SyncCron != "30 21 * * *" {
		t.Errorf("SyncCron = %q, want 30 21 * * *", cfg.SyncCron)
	}
	if cfg.FetchDelay != 1500*time.Millisecond {
		t.Errorf("FetchDelay = %v, want 1.5s", cfg.FetchDelay)
	}
	if !cfg.FetchEnable {
		t.Error("FetchEnable = false, want true")
	}
}

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("LOTTERY_ADDR", ":9090")
	t.Setenv("LOTTERY_FETCH_DELAY_MS", "500")
	t.Setenv("LOTTERY_FETCH_ENABLE", "false")

	cfg := Load()
	if cfg.Addr != ":9090" {
		t.Errorf("Addr = %q, want :9090", cfg.Addr)
	}
	if cfg.FetchDelay != 500*time.Millisecond {
		t.Errorf("FetchDelay = %v, want 500ms", cfg.FetchDelay)
	}
	if cfg.FetchEnable {
		t.Error("FetchEnable = true, want false")
	}
}
