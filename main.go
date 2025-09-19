package main

import (
	"log"
	"net/http"
	"os"

	"ismismcube-backend/internal/api"
	"ismismcube-backend/internal/config"
)

func main() {
	// 加载配置
	appConfig := config.Load()

	// 初始化 API 路由系统
	api.Init()

	// 启动服务器
	port := appConfig.Port
	if port == "" {
		port = "2998"
	}

	log.Printf("Server is running at http://127.0.0.1:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
		os.Exit(1)
	}
}
