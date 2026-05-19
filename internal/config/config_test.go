package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileSupportsEmailAndCachePathAliases(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"email":"me@example.com","password":"secret","cache_path":"/tmp/state.json"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	cfg, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile failed: %v", err)
	}
	if cfg.Username != "me@example.com" {
		t.Fatalf("Username = %q, want me@example.com", cfg.Username)
	}
	if cfg.Password != "secret" {
		t.Fatalf("Password = %q, want secret", cfg.Password)
	}
	if cfg.Cache != "/tmp/state.json" {
		t.Fatalf("Cache = %q, want /tmp/state.json", cfg.Cache)
	}
}

func TestLoadUsesEnvOverrides(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"username":"file@example.com","password":"file-pass","cache":"file-cache.json"}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("THINGS_CONFIG", path)
	t.Setenv("THINGS_USERNAME", "env@example.com")
	t.Setenv("THINGS_PASSWORD", "env-pass")
	t.Setenv("THINGS_CLI_CACHE", "env-cache.json")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Username != "env@example.com" {
		t.Fatalf("Username = %q, want env@example.com", cfg.Username)
	}
	if cfg.Password != "env-pass" {
		t.Fatalf("Password = %q, want env-pass", cfg.Password)
	}
	if cfg.Cache != "env-cache.json" {
		t.Fatalf("Cache = %q, want env-cache.json", cfg.Cache)
	}
}

func TestMissingConfigFileIsAllowed(t *testing.T) {
	cfg, err := LoadFile(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("LoadFile missing failed: %v", err)
	}
	if cfg != (Config{}) {
		t.Fatalf("cfg = %#v, want zero value", cfg)
	}
}
