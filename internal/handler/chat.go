package handler

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"ismismcube-backend/internal/manager/task_manager"
	"ismismcube-backend/internal/utility"
	"net/http"
)

type websocketCreatedResponse struct {
	WebSocketID string `json:"websocket_id"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	if len(body) == 0 {
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}
	websocketID := generateWebSocketID()
	clientIP := utility.GetRealIP(r)
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}
	task_manager.GetTaskManager().CreateChatTask(body, websocketID, clientIP, userAgent)
	utility.WriteJSON(w, http.StatusCreated, websocketCreatedResponse{WebSocketID: websocketID})
	r.Body.Close()
}

func generateWebSocketID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
