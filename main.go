package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ismismcube-backend/internal/api"
	"ismismcube-backend/internal/config"
	"ismismcube-backend/internal/manager/task_manager"
	"ismismcube-backend/internal/websocket"
)

func main() {
	config.Init()
	task_manager.InitTaskManager(&ws.WebSocketBroadcaster{})
	api.Init()

	server := &http.Server{
		Addr:    ":" + config.Port,
		Handler: nil, // 使用默认的 DefaultServeMux
	}
	go func() {
		log.Printf("Server is running at http://127.0.0.1:%s", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server shutdown completed")
	}
	config.CloseDB()
	log.Println("Database connection closed")

	log.Println("--------------------------")
	log.Println("Server exited gracefully")
}
