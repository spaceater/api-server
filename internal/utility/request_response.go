package utility

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

func GetRealIP(r *http.Request) string {
	// nginx设置X-Real-IP头已经通过frpc的transport.proxyProtocolVersion="v2"和nginx的proxy_protocol和$proxy_protocol_addr正确设置为用户的原始IP
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return strings.TrimSpace(realIP)
	}
	// 其次检查X-Forwarded-For头，通常是nginx的反向代理设置，取第一个IP并去除空格
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}
	// 如果以上都没有，则使用RemoteAddr，通常是frp的默认设置，通常是127.0.0.1
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
