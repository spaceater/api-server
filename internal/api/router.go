package api

import (
	"ismismcube-backend/internal/handler"
	"ismismcube-backend/internal/middleware"
	"ismismcube-backend/internal/router"
	"ismismcube-backend/internal/websocket"
)

func RegisterApi() {
	router.Url("/ping", handler.Ping)
	router.UrlGroup("/api",
		router.Url("/page_view", handler.PageViewHandler).Use(middleware.NoCache),
		router.Url("/send_chat", handler.ChatHandler),
	).Use(middleware.CORS)
	router.UrlGroup("/ws",
		router.Url("/ismismcube_online", ws.HandleIsmismcubeOnline),
		router.Url("/chat_broadcast", ws.HandleChatBroadcast),
		router.Url("/chat_task", ws.HandleChatTask),
	)
}

func Init() {
	RegisterApi()
	router.Init()
}
