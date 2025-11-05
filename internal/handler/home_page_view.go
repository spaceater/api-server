package handler

import (
	"ismismcube-backend/internal/server/home_server"
	"ismismcube-backend/internal/utility"
	"net/http"
)

type HomePageViewResponse struct {
	PageView int `json:"page_view"`
}

func HomePageViewHandler(w http.ResponseWriter, r *http.Request) {
	pageView, err := home_server.GetPageViewCount()
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, HomePageViewResponse{PageView: -1})
	} else {
		utility.WriteJSON(w, http.StatusOK, HomePageViewResponse{PageView: pageView + 1})
	}

	clientIP := utility.GetRealIP(r)
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}
	home_server.AddPageView(clientIP, userAgent)
}
