package ws

import (
	"ismismcube-backend/internal/config"
	"ismismcube-backend/internal/manager/task_manager"
	"log"
	"net"
	"net/http"
	"time"

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
	taskManager := task_manager.GetTaskManager()
	err = taskManager.RegisterTaskConnection(websocketID, conn)
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("task not found"))
		conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		conn.Close()
		return
	}
	conn.SetReadDeadline(time.Now().Add(config.WSPongWaitFast))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(config.WSPongWaitFast))
		return nil
	})
	ticker := time.NewTicker(config.WSPingIntervalFast)
	go func() {
		var isNormalClose bool
		defer func() {
			ticker.Stop()
			if !isNormalClose {
				if tcpConn, ok := conn.UnderlyingConn().(*net.TCPConn); ok {
					tcpConn.SetLinger(0)
				}
			}
			conn.Close()
			taskManager.UnregisterTaskConnection(websocketID)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					isNormalClose = true
				}
				break
			}
		}
	}()
	go func() {
		defer ticker.Stop()
		for {
			<-ticker.C
			conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(config.WSWriteWait))
		}
	}()
}
