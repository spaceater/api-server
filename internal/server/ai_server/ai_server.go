package ai_server

import (
	"ismismcube-backend/internal/model"
)

func GetExecutedTaskCount() (int, error) {
	return model.GetAIExecutedTaskCount()
}

func AddExecutedTask(visitorIP string) (*model.AIExecutedTask, error) {
	task := &model.AIExecutedTask{
		VisitorIP: visitorIP,
	}
	return model.AddAIExecutedTask(task)
}

func GetExecutedTasks(limit, offset int) ([]model.AIExecutedTask, error) {
	return model.GetAIExecutedTasks(limit, offset)
}

func DeleteExecutedTask(id int) error {
	return model.DeleteAIExecutedTask(id)
}
