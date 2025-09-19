package handlers

import (
	"encoding/json"
	"net/http"
)

func Ping(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"message": "pong",
		"status":  "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
