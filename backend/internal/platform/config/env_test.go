package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"erp_monolith/backend/internal/platform/config"
)

func TestEnvHelpers(t *testing.T) {
	t.Run("GetString with default and override", func(t *testing.T) {
		key := "TEST_STR_KEY"
		if val := config.GetString(key, "default_val"); val != "default_val" {
			t.Errorf("expected default_val, got %s", val)
		}

		t.Setenv(key, "custom_val")
		if val := config.GetString(key, "default_val"); val != "custom_val" {
			t.Errorf("expected custom_val, got %s", val)
		}
	})

	t.Run("GetInt and GetInt32 and GetInt64", func(t *testing.T) {
		key := "TEST_INT_KEY"
		if val := config.GetInt(key, 42); val != 42 {
			t.Errorf("expected 42, got %d", val)
		}

		t.Setenv(key, "100")
		if val := config.GetInt(key, 42); val != 100 {
			t.Errorf("expected 100, got %d", val)
		}
		if val := config.GetInt32(key, 42); val != 100 {
			t.Errorf("expected 100, got %d", val)
		}
		if val := config.GetInt64(key, 42); val != 100 {
			t.Errorf("expected 100, got %d", val)
		}

		// Invalid int should return default
		t.Setenv(key, "invalid_num")
		if val := config.GetInt(key, 42); val != 42 {
			t.Errorf("expected fallback 42, got %d", val)
		}
	})

	t.Run("GetBool", func(t *testing.T) {
		key := "TEST_BOOL_KEY"
		if val := config.GetBool(key, true); val != true {
			t.Errorf("expected true, got %v", val)
		}

		t.Setenv(key, "false")
		if val := config.GetBool(key, true); val != false {
			t.Errorf("expected false, got %v", val)
		}

		t.Setenv(key, "1")
		if val := config.GetBool(key, false); val != true {
			t.Errorf("expected true, got %v", val)
		}
	})

	t.Run("GetDuration", func(t *testing.T) {
		key := "TEST_DUR_KEY"
		if val := config.GetDuration(key, 5*time.Minute); val != 5*time.Minute {
			t.Errorf("expected 5m, got %v", val)
		}

		t.Setenv(key, "2h30m")
		if val := config.GetDuration(key, time.Minute); val != 2*time.Hour+30*time.Minute {
			t.Errorf("expected 2h30m, got %v", val)
		}
	})

	t.Run("GetStringSlice", func(t *testing.T) {
		key := "TEST_SLICE_KEY"
		def := []string{"a", "b"}
		if val := config.GetStringSlice(key, def); len(val) != 2 {
			t.Errorf("expected 2 items, got %d", len(val))
		}

		t.Setenv(key, "x, y , z")
		val := config.GetStringSlice(key, def)
		if len(val) != 3 || val[0] != "x" || val[1] != "y" || val[2] != "z" {
			t.Errorf("unexpected slice result: %v", val)
		}
	})

	t.Run("GetRequired", func(t *testing.T) {
		key := "TEST_REQ_KEY"
		if _, err := config.GetRequired(key); err == nil {
			t.Error("expected error for unset required key")
		}

		t.Setenv(key, "present")
		val, err := config.GetRequired(key)
		if err != nil || val != "present" {
			t.Errorf("expected present, got %s (err: %v)", val, err)
		}
	})
}

func TestLoadEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	content := `
# Comment line
APP_TEST_A=hello
APP_TEST_B="world with spaces"
APP_TEST_C='single quoted'
APP_TEST_D=unquoted_value # inline comment
`
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	if err := config.LoadEnvFile(envPath); err != nil {
		t.Fatalf("failed to load env file: %v", err)
	}

	if v := os.Getenv("APP_TEST_A"); v != "hello" {
		t.Errorf("expected hello, got %s", v)
	}
	if v := os.Getenv("APP_TEST_B"); v != "world with spaces" {
		t.Errorf("expected 'world with spaces', got %s", v)
	}
	if v := os.Getenv("APP_TEST_C"); v != "single quoted" {
		t.Errorf("expected 'single quoted', got %s", v)
	}
	if v := os.Getenv("APP_TEST_D"); v != "unquoted_value" {
		t.Errorf("expected 'unquoted_value', got %s", v)
	}
}
