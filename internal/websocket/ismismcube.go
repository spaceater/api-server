package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ismismcubeClients    = make(map[*websocket.Conn]struct{})
	ismismcubeClientsMux sync.RWMutex
)

func RegisterIsmismcubeClient(conn *websocket.Conn) {
	ismismcubeClientsMux.Lock()
	defer ismismcubeClientsMux.Unlock()
	ismismcubeClients[conn] = struct{}{}
	broadcastOnlineCount()
}

func UnregisterIsmismcubeClient(conn *websocket.Conn) {
	ismismcubeClientsMux.Lock()
	defer ismismcubeClientsMux.Unlock()
	if _, ok := ismismcubeClients[conn]; ok {
		delete(ismismcubeClients, conn)
		conn.Close()
		broadcastOnlineCount()
	}
}

func broadcastOnlineCount() {
	count := len(ismismcubeClients)
	data, err := json.Marshal(map[string]int{"online_count": count})
	if err != nil {
		log.Printf("JSON marshal error: %v", err)
		return
	}
	for conn := range ismismcubeClients {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Write message error: %v", err)
			delete(ismismcubeClients, conn)
			conn.Close()
		}
	}
}

func HandleIsmismcubeOnline(w http.ResponseWriter, r *http.Request) {
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
	RegisterIsmismcubeClient(conn)
	go func() {
		defer func() {
			UnregisterIsmismcubeClient(conn)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
