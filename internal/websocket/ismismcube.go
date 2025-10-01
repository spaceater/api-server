package ws

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	ismismcubeClients    = make(map[*websocket.Conn]*ClientInfo)
	ismismcubeClientsMux sync.RWMutex
)

func RegisterIsmismcubeClient(conn *websocket.Conn) {
	ismismcubeClientsMux.Lock()
	defer ismismcubeClientsMux.Unlock()
	ismismcubeClients[conn] = &ClientInfo{}
	go broadcastOnlineCount()
}

func UnregisterIsmismcubeClient(conn *websocket.Conn) {
	ismismcubeClientsMux.Lock()
	defer ismismcubeClientsMux.Unlock()
	if _, ok := ismismcubeClients[conn]; ok {
		delete(ismismcubeClients, conn)
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
		go broadcastOnlineCount()
	}
}

func broadcastOnlineCount() {
	ismismcubeClientsMux.RLock()
	data := []byte(fmt.Sprintf("broadcast:[online:%d]", len(ismismcubeClients)))
	clients := make([]*websocket.Conn, 0, len(ismismcubeClients))
	for conn := range ismismcubeClients {
		clients = append(clients, conn)
	}
	ismismcubeClientsMux.RUnlock()
	for _, conn := range clients {
		ismismcubeClientsMux.RLock()
		clientInfo, exists := ismismcubeClients[conn]
		ismismcubeClientsMux.RUnlock()
		if !exists {
			continue
		}
		clientInfo.WriteMutex.Lock()
		defer clientInfo.WriteMutex.Unlock()
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Printf("Write message error: %v", err)
			go UnregisterIsmismcubeClient(conn)
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
