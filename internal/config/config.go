package config

import (
	"os"
)

var (
	Port         string
	PageViewFile string
	ApiUrl       string
	ApiKey       string
)

func Init() {
	Port = getEnv("PORT", "2998")
	PageViewFile = getEnv("PAGE_VIEW_FILE", "./resources/page-view.txt")
	ApiUrl = getEnv("API_URL", "")
	ApiKey = getEnv("API_KEY", "")
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}
