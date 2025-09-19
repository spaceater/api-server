package api

import (
	"ismismcube-backend/internal/middleware"
	"ismismcube-backend/internal/router"
	"ismismcube-backend/internal/handlers"
	"ismismcube-backend/internal/websocket"
)

func RegisterApi() {
	router.Url("/ping", handlers.Ping)
	router.Url("/ws/ismismcube_online", websocket.HandleWebSocket)
	router.UrlGroup("/api",
		router.Url("/page_view", handlers.PageViewHandler).Use(middleware.NoCache),
	)
}

func Init() {
	RegisterApi()
	router.Init()
}
