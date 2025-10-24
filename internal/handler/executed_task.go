package handler

import (
	"ismismcube-backend/internal/server/ai_server"
	"net/http"
)

type ExecutedTaskResponse struct {
	ExecutedTask int `json:"executed_task"`
}

func ExecutedTaskHandler(w http.ResponseWriter, r *http.Request) {
	executedTask, err := ai_server.GetExecutedTaskCount()
	if err != nil {
		sendResponse(w, ExecutedTaskResponse{ExecutedTask: -1})
		return
	}
	sendResponse(w, ExecutedTaskResponse{ExecutedTask: executedTask})
}
