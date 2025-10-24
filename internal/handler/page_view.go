package handler

import (
	"ismismcube-backend/internal/server/ismismcube_server"
	"ismismcube-backend/internal/utility"
	"net/http"
)

type PageViewResponse struct {
	PageView int `json:"page_view"`
}

func PageViewHandler(w http.ResponseWriter, r *http.Request) {
	pageView, err := ismismcube_server.GetPageViewCount()
	if err != nil {
		utility.WriteJSON(w, http.StatusInternalServerError, PageViewResponse{PageView: -1})
	} else {
		utility.WriteJSON(w, http.StatusOK, PageViewResponse{PageView: pageView + 1})
	}

	clientIP := utility.GetRealIP(r)
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}
	ismismcube_server.AddPageView(clientIP, userAgent)
}
