package handler

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

var pageViewMutex sync.Mutex

func PageViewHandler(w http.ResponseWriter, r *http.Request) {
	pageViewMutex.Lock()
	defer pageViewMutex.Unlock()
	pageViewFile := config.PageViewFile
	data, err := os.ReadFile(pageViewFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建文件并初始化为1
			if err := os.WriteFile(pageViewFile, []byte("1"), 0644); err != nil {
				sendResponse(w, PageViewResponse{PageView: -1})
				return
			}
			sendResponse(w, PageViewResponse{PageView: 1})
			return
		}
		sendResponse(w, PageViewResponse{PageView: -1})
		return
	}
	pageView, err := strconv.Atoi(string(data))
	if err != nil {
		// 文件内容不是有效数字，重置为1
		if err := os.WriteFile(pageViewFile, []byte("1"), 0644); err != nil {
			sendResponse(w, PageViewResponse{PageView: -1})
			return
		}
		sendResponse(w, PageViewResponse{PageView: 1})
		return
	}
	pageView++
	if err := os.WriteFile(pageViewFile, []byte(strconv.Itoa(pageView)), 0644); err != nil {
		sendResponse(w, PageViewResponse{PageView: -1})
		return
	}
	sendResponse(w, PageViewResponse{PageView: pageView})
}

type ExecutedTaskResponse struct {
	ExecutedTask int `json:"executed_task"`
}

var executedTaskMutex sync.Mutex

func ExecutedTaskHandler(w http.ResponseWriter, r *http.Request) {
	executedTaskMutex.Lock()
	defer executedTaskMutex.Unlock()
	executedTaskFile := config.ExecutedTaskFile
	data, err := os.ReadFile(executedTaskFile)
	if err != nil {
		if os.IsNotExist(err) {
			// 文件不存在，创建文件并初始化为1
			if err := os.WriteFile(executedTaskFile, []byte("1"), 0644); err != nil {
				sendResponse(w, ExecutedTaskResponse{ExecutedTask: -1})
				return
			}
			sendResponse(w, ExecutedTaskResponse{ExecutedTask: 1})
			return
		}
		sendResponse(w, ExecutedTaskResponse{ExecutedTask: -1})
		return
	}
	executedTask, err := strconv.Atoi(string(data))
	if err != nil {
		// 文件内容不是有效数字，重置为1
		if err := os.WriteFile(executedTaskFile, []byte("1"), 0644); err != nil {
			sendResponse(w, ExecutedTaskResponse{ExecutedTask: -1})
			return
		}
		sendResponse(w, ExecutedTaskResponse{ExecutedTask: 1})
		return
	}
	if err := os.WriteFile(executedTaskFile, []byte(strconv.Itoa(executedTask)), 0644); err != nil {
		sendResponse(w, ExecutedTaskResponse{ExecutedTask: -1})
		return
	}
	sendResponse(w, ExecutedTaskResponse{ExecutedTask: executedTask})
}

type WebStatisticsHandler struct{}

func (wsh *WebStatisticsHandler) IncrementExecutedTask() error {
	executedTaskMutex.Lock()
	defer executedTaskMutex.Unlock()
	executedTaskFile := config.ExecutedTaskFile
	data, err := os.ReadFile(executedTaskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(executedTaskFile, []byte("1"), 0644)
		}
		return err
	}
	executedTask, err := strconv.Atoi(string(data))
	if err != nil {
		return os.WriteFile(executedTaskFile, []byte("1"), 0644)
	}
	executedTask++
	return os.WriteFile(executedTaskFile, []byte(strconv.Itoa(executedTask)), 0644)
}

func sendResponse(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
