package ws

import (
	"encoding/json"
	"ismismcube-backend/internal/server"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	chatClients    = make(map[*websocket.Conn]struct{})
	chatClientsMux sync.RWMutex
)

type WebSocketBroadcaster struct{}

func (w *WebSocketBroadcaster) BroadcastQueueStats(waiting, executing int) {
	chatClientsMux.Lock()
	defer chatClientsMux.Unlock()
	broadcastQueueStats(waiting, executing)
}

func RegisterBroadcastClient(conn *websocket.Conn, waiting, executing int) {
	chatClientsMux.Lock()
	defer chatClientsMux.Unlock()
	chatClients[conn] = struct{}{}
	data, err := json.Marshal(map[string]interface{}{
		"waiting_count":   waiting,
		"executing_count": executing,
	})
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}
	sendQueueStatus(conn, data)
}

func UnregisterBroadcastClient(conn *websocket.Conn) {
	chatClientsMux.Lock()
	defer chatClientsMux.Unlock()
	if _, ok := chatClients[conn]; ok {
		delete(chatClients, conn)
		conn.Close()
	}
}

// 调用此函数前需要确保chatClientsMux已经锁定
func sendQueueStatus(conn *websocket.Conn, data []byte) {
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Write message error: %v", err)
		delete(chatClients, conn)
		conn.Close()
	}
}

// 调用此函数前需要确保chatClientsMux已经锁定
func broadcastQueueStats(waiting, executing int) {
	data, err := json.Marshal(map[string]interface{}{
		"waiting_count":   waiting,
		"executing_count": executing,
	})
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}
	clients := make([]*websocket.Conn, 0, len(chatClients))
	for conn := range chatClients {
		clients = append(clients, conn)
	}
	for _, conn := range clients {
		sendQueueStatus(conn, data)
	}
}

func HandleChatBroadcast(w http.ResponseWriter, r *http.Request) {
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
	waiting, executing := server.GetTaskManager().GetQueueCount()
	RegisterBroadcastClient(conn, waiting, executing)
	go func() {
		defer func() {
			UnregisterBroadcastClient(conn)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
