package ws

import (
	"ismismcube-backend/internal/config"
	"ismismcube-backend/internal/utility"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type OnlineCountData struct {
	Online int `json:"online"`
}

type DanmuData struct {
	Content string `json:"content"`
}

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
	delete(ismismcubeClients, conn)
	go broadcastOnlineCount()
}

func broadcastOnlineCount() {
	ismismcubeClientsMux.RLock()
	data := &utility.MessageData{
		Type: "broadcast",
		Data: OnlineCountData{
			Online: len(ismismcubeClients),
		},
	}
	clients := make([]*websocket.Conn, 0, len(ismismcubeClients))
	for conn := range ismismcubeClients {
		clients = append(clients, conn)
	}
	ismismcubeClientsMux.RUnlock()
	msg, err := data.ToBytes()
	if err != nil {
		return
	}
	for _, conn := range clients {
		ismismcubeClientsMux.RLock()
		clientInfo, exists := ismismcubeClients[conn]
		ismismcubeClientsMux.RUnlock()
		if !exists {
			continue
		}
		clientInfo.WriteMutex.Lock()
		conn.WriteMessage(websocket.TextMessage, msg)
		clientInfo.WriteMutex.Unlock()
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
	go BroadcastIsmismcubeDanmu("有新的浏览者打开了此网页")
	conn.SetReadDeadline(time.Now().Add(config.WSPongWaitSlow))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(config.WSPongWaitSlow))
		return nil
	})

	ticker := time.NewTicker(config.WSPingIntervalSlow)
	done := make(chan struct{})
	go func() {
		var isNormalClose bool
		defer func() {
			close(done)
			ticker.Stop()
			if !isNormalClose {
				if tcpConn, ok := conn.UnderlyingConn().(*net.TCPConn); ok {
					tcpConn.SetLinger(0)
				}
			}
			conn.Close()
			UnregisterIsmismcubeClient(conn)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					isNormalClose = true
				}
				return
			}
		}
	}()
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				select {
				case <-done:
					return
				default:
					if err := conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(config.WSWriteWait)); err != nil {
						return
					}
				}
			}
		}
	}()
}

func BroadcastIsmismcubeDanmu(content string) {
	ismismcubeClientsMux.RLock()
	clients := make([]*websocket.Conn, 0, len(ismismcubeClients))
	for conn := range ismismcubeClients {
		clients = append(clients, conn)
	}
	ismismcubeClientsMux.RUnlock()
	data := &utility.MessageData{
		Type: "danmu",
		Data: DanmuData{
			Content: content,
		},
	}
	msg, err := data.ToBytes()
	if err != nil {
		return
	}
	for _, conn := range clients {
		ismismcubeClientsMux.RLock()
		clientInfo, exists := ismismcubeClients[conn]
		ismismcubeClientsMux.RUnlock()
		if !exists {
			continue
		}
		clientInfo.WriteMutex.Lock()
		conn.WriteMessage(websocket.TextMessage, msg)
		clientInfo.WriteMutex.Unlock()
	}
}
