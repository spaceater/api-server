package handler

import (
	"encoding/json"
	"io"
	"ismismcube-backend/internal/server/ismismcube_server"
	"ismismcube-backend/internal/utility"
	"net/http"
)

type SendDanmuRequest struct {
	Content string `json:"content"`
}

type SendDanmuResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func SendDanmuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, SendDanmuResponse{
			Success: false,
			Message: "Failed to read request body",
		})
		return
	}
	defer r.Body.Close()

	if len(body) == 0 {
		utility.WriteJSON(w, http.StatusBadRequest, SendDanmuResponse{
			Success: false,
			Message: "Request body is empty",
		})
		return
	}

	var req SendDanmuRequest
	if err := json.Unmarshal(body, &req); err != nil {
		utility.WriteJSON(w, http.StatusBadRequest, SendDanmuResponse{
			Success: false,
			Message: "Invalid JSON format",
		})
		return
	}

	if req.Content == "" {
		utility.WriteJSON(w, http.StatusBadRequest, SendDanmuResponse{
			Success: false,
			Message: "Content cannot be empty",
		})
		return
	}

	clientIP := utility.GetRealIP(r)
	err = ismismcube_server.SendDanmu(clientIP, req.Content)
	if err != nil {
		utility.WriteJSON(w, http.StatusTooManyRequests, SendDanmuResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	utility.WriteJSON(w, http.StatusOK, SendDanmuResponse{
		Success: true,
		Message: "Danmu sent successfully",
	})
}
