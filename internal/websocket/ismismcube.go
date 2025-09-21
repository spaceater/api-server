package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 全局变量
var (
	clients    = make(map[*websocket.Conn]struct{})
	clientsMux sync.RWMutex
)

func RegisterClient(conn *websocket.Conn) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	clients[conn] = struct{}{}
	broadcastOnlineCountLocked()
}

func UnregisterClient(conn *websocket.Conn) {
	clientsMux.Lock()
	defer clientsMux.Unlock()
	if _, ok := clients[conn]; ok {
		delete(clients, conn)
		conn.Close()
		broadcastOnlineCountLocked()
	}
}

// broadcastOnlineCountLocked 广播在线用户数量（假设已持有锁）
func broadcastOnlineCountLocked() {
	count := len(clients)
	// 直接使用map构建JSON消息
	data, err := json.Marshal(map[string]int{"online_count": count})
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}

	// 直接在锁内广播，避免并发写入
	for conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Write message error: %v", err)
			// 写入失败的连接将被标记为需要清理
			delete(clients, conn)
			conn.Close()
		}
	}
}

// HandleWebSocket 处理 WebSocket 连接
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
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

	// 注册连接
	RegisterClient(conn)

	// 处理连接关闭
	go func() {
		defer func() {
			UnregisterClient(conn)
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
