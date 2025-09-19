package config

import (
	"os"
)

// Config 应用配置结构
type Config struct {
	Port         string
	PageViewFile string
}

// Load 加载配置
func Load() *Config {
	return &Config{
		Port:         getEnv("PORT", "2998"),
		PageViewFile: getEnv("PAGE_VIEW_FILE", "resources/page-view.txt"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
