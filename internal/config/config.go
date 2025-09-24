package config

import (
	"os"
	"strconv"
)

var (
	Port               string
	PageViewFile       string
	ApiUrl             string
	ApiKey             string
	MaxConcurrentTasks int
)

func Init() {
	Port = getEnv("PORT", "2998")
	PageViewFile = getEnv("PAGE_VIEW_FILE", "./resources/page-view.txt")
	ApiUrl = getEnv("API_URL", "")
	ApiKey = getEnv("API_KEY", "")
	MaxConcurrentTasks = getEnvInt("MAX_CONCURRENT_TASKS", 4)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
