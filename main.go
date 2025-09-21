package main

import (
	"log"
	"net/http"
	"os"

	"ismismcube-backend/internal/api"
	"ismismcube-backend/internal/config"
)

func main() {
	config.Init()

	api.Init()

	log.Printf("Server is running at http://127.0.0.1:%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal("Failed to start server:", err)
		os.Exit(1)
	}
}
