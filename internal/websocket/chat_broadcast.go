package ws

import (
	"encoding/json"
	"fmt"
	"ismismcube-backend/internal/config"
	"ismismcube-backend/internal/server"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type ClientInfo struct {
	WriteMutex sync.Mutex
}

var (
	chatClients    = make(map[*websocket.Conn]*ClientInfo)
	chatClientsMux sync.RWMutex
)

type WebSocketBroadcaster struct{}

func (w *WebSocketBroadcaster) BroadcastQueueStats(waiting, executing int) {
	broadcastQueueStats(waiting, executing)
}

func RegisterChatClient(conn *websocket.Conn, waiting, executing int) {
	chatClientsMux.Lock()
	defer chatClientsMux.Unlock()
	chatClients[conn] = &ClientInfo{}
	go sendQueueStats(conn, []byte(fmt.Sprintf("broadcast:[waiting_count:%d,executing_count:%d]", waiting, executing)))
	message := map[string]interface{}{
		"max_concurrent_tasks": config.LLMConfigure.MaxConcurrentTasks,
		"available_models":     config.LLMConfigure.AvailableModels,
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		return
	}
	data := []byte(fmt.Sprintf("server-config:%s", string(jsonData)))
	go sendQueueStats(conn, data)
}

func broadcastQueueStats(waiting, executing int) {
	chatClientsMux.RLock()
	clients := make([]*websocket.Conn, 0, len(chatClients))
	for conn := range chatClients {
		clients = append(clients, conn)
	}
	chatClientsMux.RUnlock()
	data := []byte(fmt.Sprintf("broadcast:[waiting_count:%d,executing_count:%d]", waiting, executing))
	for _, conn := range clients {
		sendQueueStats(conn, data)
	}
}

func sendQueueStats(conn *websocket.Conn, data []byte) {
	chatClientsMux.RLock()
	clientInfo, exists := chatClients[conn]
	chatClientsMux.RUnlock()
	if !exists {
		return
	}
	clientInfo.WriteMutex.Lock()
	defer clientInfo.WriteMutex.Unlock()
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Write message error: %v", err)
		go UnregisterChatClient(conn)
	}
}

func UnregisterChatClient(conn *websocket.Conn) {
	chatClientsMux.Lock()
	defer chatClientsMux.Unlock()
	if _, ok := chatClients[conn]; ok {
		delete(chatClients, conn)
		conn.Close()
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
	RegisterChatClient(conn, waiting, executing)
	go func() {
		defer func() {
			UnregisterChatClient(conn)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
