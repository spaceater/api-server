package handlers

import (
	"encoding/json"
	"ismismcube-backend/internal/config"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type PageViewResponse struct {
	PageView int `json:"page_view"`
}

// pageViewMutex 保护页面访问计数的互斥锁
var pageViewMutex sync.Mutex

func PageViewHandler(w http.ResponseWriter, r *http.Request) {
	// 使用互斥锁确保文件操作的原子性，避免并发竞态条件
	pageViewMutex.Lock()
	defer pageViewMutex.Unlock()

	pageViewFile := config.PageViewFile

	// 读取当前页面访问计数
	data, err := os.ReadFile(pageViewFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建文件并初始化为1
			if err := os.WriteFile(pageViewFile, []byte("1"), 0644); err != nil {
				sendErrorResponse(w, -1)
				return
			}
			sendSuccessResponse(w, 1)
			return
		}
		// 其他读取错误
		sendErrorResponse(w, -1)
		return
	}

	// 解析当前计数值
	pageView, err := strconv.Atoi(string(data))
	if err != nil {
		// 文件内容不是有效数字，重置为1
		if err := os.WriteFile(pageViewFile, []byte("1"), 0644); err != nil {
			sendErrorResponse(w, -1)
			return
		}
		sendSuccessResponse(w, 1)
		return
	}

	// 增加计数并写回文件
	pageView++
	if err := os.WriteFile(pageViewFile, []byte(strconv.Itoa(pageView)), 0644); err != nil {
		sendErrorResponse(w, -1)
		return
	}

	sendSuccessResponse(w, pageView)
}

// sendSuccessResponse 发送成功响应
func sendSuccessResponse(w http.ResponseWriter, pageView int) {
	response := PageViewResponse{PageView: pageView}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// sendErrorResponse 发送错误响应
func sendErrorResponse(w http.ResponseWriter, pageView int) {
	response := PageViewResponse{PageView: pageView}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
