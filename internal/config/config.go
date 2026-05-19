package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const envConfigPath = "THINGS_CONFIG"

type Config struct {
	Username string
	Password string
	Cache    string
}

type fileConfig struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cache     string `json:"cache"`
	CachePath string `json:"cache_path"`
}

func Path() string {
	if path := os.Getenv(envConfigPath); path != "" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ""
	}
	return filepath.Join(home, ".things-cloud.json")
}

func Load() (Config, error) {
	cfg := Config{}
	path := Path()
	if path != "" {
		fileCfg, err := LoadFile(path)
		if err != nil {
			return Config{}, err
		}
		cfg = fileCfg
	}
	return ApplyEnv(cfg), nil
}

func LoadFile(path string) (Config, error) {
	bs, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}

	var raw fileConfig
	if err := json.Unmarshal(bs, &raw); err != nil {
		return Config{}, fmt.Errorf("parse config %s: %w", path, err)
	}
	cfg := Config{
		Username: raw.Username,
		Password: raw.Password,
		Cache:    raw.Cache,
	}
	if cfg.Username == "" {
		cfg.Username = raw.Email
	}
	if cfg.Cache == "" {
		cfg.Cache = raw.CachePath
	}
	return cfg, nil
}

func ApplyEnv(cfg Config) Config {
	if v := os.Getenv("THINGS_USERNAME"); v != "" {
		cfg.Username = v
	}
	if v := os.Getenv("THINGS_PASSWORD"); v != "" {
		cfg.Password = v
	}
	if v := os.Getenv("THINGS_CLI_CACHE"); v != "" {
		cfg.Cache = v
	}
	return cfg
}
