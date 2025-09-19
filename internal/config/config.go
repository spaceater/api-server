package config

import (
	"os"
	"path/filepath"
)

// Config 应用配置结构
type Config struct {
	Port         string
	PageViewFile string
}

// Load 加载配置
func Load() *Config {
	pageViewFile := os.Getenv("PAGE_VIEW_FILE")
	if pageViewFile == "" {
		// 获取可执行文件所在目录
		execPath, err := os.Executable()
		if err != nil {
			// 如果无法获取可执行文件路径，使用当前工作目录
			execPath = "."
		}
		execDir := filepath.Dir(execPath)
		// 检查可执行文件名是否为 "server"（生产环境）
		if filepath.Base(execPath) == "server" {
			// 生产环境：使用 backend/resources/page-view.txt
			pageViewFile = filepath.Join(execDir, "resources", "page-view.txt")
		} else {
			// 开发环境：使用 ./resources/page-view.txt
			pageViewFile = "./resources/page-view.txt"
		}
	}

	return &Config{
		Port:         getEnv("PORT", "2998"),
		PageViewFile: pageViewFile,
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
