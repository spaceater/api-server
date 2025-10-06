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

func (w *WebSocketBroadcaster) BroadcastQueueStats(waiting, executing int, broadcastFlag int64) {
	broadcastQueueStats(waiting, executing, broadcastFlag)
}

func RegisterChatClient(conn *websocket.Conn, waiting, executing int, broadcastFlag int64) {
	chatClientsMux.Lock()
	defer chatClientsMux.Unlock()
	chatClients[conn] = &ClientInfo{}
	go sendQueueStats(conn, []byte(fmt.Sprintf(`broadcast:{"waiting_count":%d,"executing_count":%d,"broadcast_flag":%d}`, waiting, executing, broadcastFlag)))
  llmConfigData, err := json.Marshal(map[string]interface{}{
		"max_concurrent_tasks": config.LLMConfigure.MaxConcurrentTasks,
		"available_models":     config.LLMConfigure.AvailableModels,
	})
	if err != nil {
		return
	}
	go sendQueueStats(conn, []byte(fmt.Sprintf("server-config:%s", string(llmConfigData))))
	chatParamsData, err := json.Marshal(config.ChatParameters)
  if err != nil {
		return
	}
	go sendQueueStats(conn, []byte(fmt.Sprintf("chat-config:%s", string(chatParamsData))))
}

func broadcastQueueStats(waiting, executing int, broadcastFlag int64) {
	chatClientsMux.RLock()
	clients := make([]*websocket.Conn, 0, len(chatClients))
	for conn := range chatClients {
		clients = append(clients, conn)
	}
	chatClientsMux.RUnlock()
	data := []byte(fmt.Sprintf(`broadcast:{"waiting_count":%d,"executing_count":%d,"broadcast_flag":%d}`, waiting, executing, broadcastFlag))
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
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
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
	broadcastFlag := server.GetTaskManager().GetBroadcastFlag()
	RegisterChatClient(conn, waiting, executing, broadcastFlag)
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
