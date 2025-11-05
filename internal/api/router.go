package api

import (
	"ismismcube-backend/internal/handler"
	"ismismcube-backend/internal/middleware"
	"ismismcube-backend/internal/router"
  "ismismcube-backend/internal/websocket"
)

func RegisterApi() {
	router.Url("/ping", handler.Ping)
	router.UrlGroup("/home",
		router.Url("/page_view", handler.HomePageViewHandler).Use(middleware.NoCache),
	)
	router.UrlGroup("/ismismcube",
		router.Url("/page_view", handler.IsmismcubePageViewHandler).Use(middleware.NoCache),
		router.Url("/online", ws.HandleIsmismcubeOnline),
	)
	router.UrlGroup("/ai",
		router.Url("/executed_task", handler.ExecutedTaskHandler).Use(middleware.NoCache),
		router.Url("/send_chat", handler.ChatHandler),
		router.Url("/chat_broadcast", ws.HandleChatBroadcast),
		router.Url("/chat_task", ws.HandleChatTask),
	)
}

func Init() {
	RegisterApi()
	router.Init()
}
