package handler

import (
	"ismismcube-backend/internal/server/ismismcube_server"
	"net"
	"net/http"
)

type PageViewResponse struct {
	PageView int `json:"page_view"`
}

func PageViewHandler(w http.ResponseWriter, r *http.Request) {
	pageView, err := ismismcube_server.GetPageViewCount()
	if err != nil {
		sendResponse(w, PageViewResponse{PageView: -1})
	} else {
		sendResponse(w, PageViewResponse{PageView: pageView + 1})
	}

	var clientIP string
	clientIP, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		clientIP = r.RemoteAddr
	}
	var userAgent string
	userAgent = r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}
	ismismcube_server.AddPageView(clientIP, userAgent)
}
