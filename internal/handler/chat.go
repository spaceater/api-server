package handler

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"ismismcube-backend/internal/manager/task_manager"
	"net/http"
)

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
	task_manager.GetTaskManager().CreateChatTask(body, websocketID)
	response := map[string]string{
		"websocket_id": websocketID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	r.Body.Close()
}

func generateWebSocketID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
