package handler

import (
	"ismismcube-backend/internal/utility"
  "net/http"
)

type pingResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
func Ping(w http.ResponseWriter, r *http.Request) {
	utility.WriteJSON(w, http.StatusOK, pingResponse{Message: "pong", Status: "ok"})
}
