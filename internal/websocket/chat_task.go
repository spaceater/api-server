package ws

import (
	"ismismcube-backend/internal/server"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func HandleChatTask(w http.ResponseWriter, r *http.Request) {
	websocketID := r.URL.Query().Get("id")
	if websocketID == "" {
		http.Error(w, "Missing websocket ID", http.StatusBadRequest)
		return
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // 允许所有来源
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	taskManager := server.GetTaskManager()
	taskManager.RegisterTaskConnection(websocketID, conn)
	go func() {
		defer taskManager.UnregisterTaskConnection(websocketID)
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
