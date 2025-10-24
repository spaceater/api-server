package handler

import (
	"ismismcube-backend/internal/server/ai_server"
	"ismismcube-backend/internal/utility"
	"net/http"
)

type ExecutedTaskResponse struct {
	ExecutedTask int `json:"executed_task"`
}

func ExecutedTaskHandler(w http.ResponseWriter, r *http.Request) {
	executedTask, err := ai_server.GetExecutedTaskCount()
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, ExecutedTaskResponse{ExecutedTask: -1})
		return
	}
	utility.WriteJSON(w, http.StatusOK, ExecutedTaskResponse{ExecutedTask: executedTask})
}
