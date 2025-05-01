package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all the application configuration
type Config struct {
	Server  ServerConfig  `toml:"server"`
	Content ContentConfig `toml:"content"`
	Redis   RedisConfig   `toml:"redis"`
	I18n    I18nConfig    `toml:"i18n"`
	Theme   ThemeConfig   `toml:"theme"`
}

// ServerConfig holds server-related settings
type ServerConfig struct {
	Port    string `toml:"port"`
	BaseURL string `toml:"base_url"`
}

// ContentConfig holds content-related settings
type ContentConfig struct {
	PostsDir string `toml:"posts_dir"`
	SiteName string `toml:"site_name"`
}

// RedisConfig holds Redis-related settings
type RedisConfig struct {
	URL      string        `toml:"url"`
	CacheTTL time.Duration `toml:"cache_ttl"`
}

// I18nConfig holds i18n settings
type I18nConfig struct {
	TranslationsDir string `toml:"translations_dir"`
	DefaultLang     string `toml:"default_lang"`
}

// ThemeConfig holds theme settings
type ThemeConfig struct {
	DefaultTheme string `toml:"default_theme"` // "light" or "dark"
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    ":8080",
			BaseURL: "http://localhost:8080",
		},
		Content: ContentConfig{
			PostsDir: "./posts",
			SiteName: "My Blog",
		},
		Redis: RedisConfig{
			URL:      "localhost:6379",
			CacheTTL: 10 * time.Minute,
		},
		I18n: I18nConfig{
			TranslationsDir: "./translations",
			DefaultLang:     "en",
		},
		Theme: ThemeConfig{
			DefaultTheme: "light",
		},
	}
}

// LoadConfig loads the application configuration from a TOML file
func LoadConfig(configPath string) (*Config, error) {
	// Use default config if no path specified
	if configPath == "" {
		configPath = "config.toml"
	}

	// Check if config file exists, return default if not
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	// Read and parse the config file
	config := DefaultConfig()
	if _, err := toml.DecodeFile(configPath, config); err != nil {
		return nil, err
	}

	// Ensure posts directory is an absolute path
	if !filepath.IsAbs(config.Content.PostsDir) {
		absPath, err := filepath.Abs(config.Content.PostsDir)
		if err == nil {
			config.Content.PostsDir = absPath
		}
	}

	// Ensure translations directory is an absolute path
	if !filepath.IsAbs(config.I18n.TranslationsDir) {
		absPath, err := filepath.Abs(config.I18n.TranslationsDir)
		if err == nil {
			config.I18n.TranslationsDir = absPath
		}
	}

	return config, nil
}
