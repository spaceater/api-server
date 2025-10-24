package model

import (
	"ismismcube-backend/internal/config"
	"log"
	"time"
)

type AIExecutedTask struct {
	ID        int       `json:"id"`
	Time      time.Time `json:"time"`
	VisitorIP string    `json:"visitor_ip"`
}

func GetAIExecutedTaskCount() (int, error) {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM ai_executed_tasks").Scan(&count)
	if err != nil {
		log.Printf("Error getting AI executed task count: %v", err)
		return -1, err
	}
	return count, nil
}

func AddAIExecutedTask(task *AIExecutedTask) (*AIExecutedTask, error) {
	result, err := config.DB.Exec(
		"INSERT INTO ai_executed_tasks (visitor_ip) VALUES (?)",
		task.VisitorIP,
	)
	if err != nil {
		log.Printf("Error adding AI executed task: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		return nil, err
	}

	task.ID = int(id)
	return task, nil
}

func GetAIExecutedTasks(limit, offset int) ([]AIExecutedTask, error) {
	rows, err := config.DB.Query(
		"SELECT id, time, visitor_ip FROM ai_executed_tasks ORDER BY time DESC LIMIT ? OFFSET ?",
		limit, offset,
	)
	if err != nil {
		log.Printf("Error getting AI executed tasks: %v", err)
		return nil, err
	}
	defer rows.Close()

	var tasks []AIExecutedTask
	for rows.Next() {
		var task AIExecutedTask
		err := rows.Scan(&task.ID, &task.Time, &task.VisitorIP)
		if err != nil {
			log.Printf("Error scanning AI executed task: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func DeleteAIExecutedTask(id int) error {
	_, err := config.DB.Exec("DELETE FROM ai_executed_tasks WHERE id = ?", id)
	if err != nil {
		log.Printf("Error deleting AI executed task: %v", err)
		return err
	}

	return nil
}
