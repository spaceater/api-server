package server

import (
	"bytes"
	"log"
	"fmt"
	"io"
	"ismismcube-backend/internal/config"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ChatTask struct {
	CreatedAt     time.Time       `json:"created_at"`
	Content       []byte          `json:"-"`
	WebSocketID   string          `json:"websocket_id"`
	WebSocketConn *websocket.Conn `json:"-"`
	WriteMutex         sync.Mutex      `json:"-"`
}

type QueueBroadcaster interface {
	BroadcastQueueStats(waiting, executing int)
}

type TaskManager struct {
	pendingTasks   map[string]*ChatTask
	waitingTasks   []*ChatTask
	executingTasks map[string]*ChatTask
	mutex          sync.RWMutex
	broadcaster    QueueBroadcaster
}

var (
	taskManager *TaskManager
)

func GetTaskManager() *TaskManager {
	return taskManager
}

func InitTaskManager(broadcaster QueueBroadcaster) {
	taskManager = &TaskManager{
		pendingTasks:   make(map[string]*ChatTask),
		waitingTasks:   make([]*ChatTask, 0),
		executingTasks: make(map[string]*ChatTask),
	}
	taskManager.broadcaster = broadcaster
}

func (tm *TaskManager) CreateChatTask(content []byte, websocketID string) {
	task := &ChatTask{
		CreatedAt:   time.Now(),
		Content:     content,
		WebSocketID: websocketID,
	}
	tm.mutex.Lock()
	tm.pendingTasks[websocketID] = task
	tm.mutex.Unlock()
	// 必须创建一个新的goroutine，防止CreateChatTask被阻塞
	go func() {
		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()
		<-timer.C
		tm.mutex.Lock()
		delete(tm.pendingTasks, websocketID)
		tm.mutex.Unlock()
	}()
}

func (tm *TaskManager) RegisterTaskConnection(websocketID string, conn *websocket.Conn) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	if task, exists := tm.pendingTasks[websocketID]; exists {
		task.WebSocketConn = conn
		tm.waitingTasks = append(tm.waitingTasks, task)
		delete(tm.pendingTasks, websocketID)
		go tm.broadcaster.BroadcastQueueStats(len(tm.waitingTasks), len(tm.executingTasks))
		go tm.sendTaskPosition(task, len(tm.waitingTasks)-1)
		// 触发任务调度
		go tm.checkTasks()
		return
	}
	// 执行中的任务允许重连
	if task, exists := tm.executingTasks[websocketID]; exists {
		task.WebSocketConn = conn
		return
	}
}

func (tm *TaskManager) UnregisterTaskConnection(websocketID string) {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	for i, task := range tm.waitingTasks {
		if task.WebSocketID == websocketID {
			if task.WebSocketConn != nil {
				task.WebSocketConn.Close()
				task.WebSocketConn = nil
			}
			tm.waitingTasks = append(tm.waitingTasks[:i], tm.waitingTasks[i+1:]...)
			go tm.broadcaster.BroadcastQueueStats(len(tm.waitingTasks), len(tm.executingTasks))
			go tm.broadcastTasksPositions()
			return
		}
	}
	// 执行中的任务断开后保留在executingTasks中，留给callLLM处理
	if task, exists := tm.executingTasks[websocketID]; exists {
		if task.WebSocketConn != nil {
			task.WebSocketConn.Close()
			task.WebSocketConn = nil
		}
		return
	}
}

func (tm *TaskManager) checkTasks() {
	tm.mutex.Lock()
	defer tm.mutex.Unlock()
	if len(tm.executingTasks) >= config.LLMConfigure.MaxConcurrentTasks {
		return
	}
	tasksScheduled := false
	for i := 0; i < len(tm.waitingTasks); i++ {
		if len(tm.executingTasks) >= config.LLMConfigure.MaxConcurrentTasks {
			break
		}
		tasksScheduled = true
		task := tm.waitingTasks[i]
		tm.waitingTasks = append(tm.waitingTasks[:i], tm.waitingTasks[i+1:]...)
		i--
		if task.WebSocketConn != nil {
			tm.executingTasks[task.WebSocketID] = task
			go tm.executeTask(task)
			go tm.sendTaskPosition(task, -1)
		}
	}
	if tasksScheduled {
		go tm.broadcastTasksPositions()
		go tm.broadcaster.BroadcastQueueStats(len(tm.waitingTasks), len(tm.executingTasks))
	}
}

func (tm *TaskManager) executeTask(task *ChatTask) {
	defer func() {
		if task.WebSocketConn != nil {
			task.WebSocketConn.Close()
		}
		tm.mutex.Lock()
		delete(tm.executingTasks, task.WebSocketID)
		tm.mutex.Unlock()
		go tm.broadcaster.BroadcastQueueStats(tm.GetQueueCount())
		go tm.checkTasks()
	}()
	tm.callLLM(task)
}

func (tm *TaskManager) callLLM(task *ChatTask) {
	tm.mutex.RLock()
	conn := task.WebSocketConn
	tm.mutex.RUnlock()
	if conn == nil {
		return
	}
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	req, err := http.NewRequest("POST", config.LLMConfigure.ApiUrl, bytes.NewBuffer(task.Content))
	if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("data: {\"error\": \"Failed to create request\"}\n\n"))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Authorization", "Bearer "+config.LLMConfigure.ApiKey)
	resp, err := client.Do(req)
	if err != nil {
    log.Println("Failed to send request to AI API", err)
		conn.WriteMessage(websocket.TextMessage, []byte("data: {\"error\": \"Failed to send request to AI API\"}\n\n"))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		errorBody, _ := io.ReadAll(resp.Body)
		errorMsg := fmt.Sprintf("data: {\"error\": \"AI API returned status %d: %s\"}\n\n", resp.StatusCode, string(errorBody))
		conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
		return
	}
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		tm.mutex.RLock()
		conn := task.WebSocketConn
		tm.mutex.RUnlock()
		if conn == nil {
			break
		}
		if n > 0 {
			if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}
}

func (tm *TaskManager) GetQueueCount() (waiting, executing int) {
	tm.mutex.RLock()
	waiting = len(tm.waitingTasks)
	executing = len(tm.executingTasks)
	tm.mutex.RUnlock()
	return waiting, executing
}

func (tm *TaskManager) broadcastTasksPositions() {
	tm.mutex.RLock()
	tasks := make([]*ChatTask, len(tm.waitingTasks))
	copy(tasks, tm.waitingTasks)
	tm.mutex.RUnlock()
	for i, task := range tasks {
		tm.sendTaskPosition(task, i)
	}
}

func (tm *TaskManager) sendTaskPosition(task *ChatTask, position int) {
	conn := task.WebSocketConn
	if conn == nil {
		go tm.UnregisterTaskConnection(task.WebSocketID)
		return
	}
	data := []byte(fmt.Sprintf("broadcast:[queue_position:%d]", position))
	// websocket连接只允许同时传输一个消息，所以需要锁定
	task.WriteMutex.Lock()
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		go tm.UnregisterTaskConnection(task.WebSocketID)
	}
	task.WriteMutex.Unlock()
}
