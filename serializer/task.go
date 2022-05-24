package serializer

import (
	"gin-memos/model"
)

type TaskResponse struct {
	// gorm 用 tag 的方式来标识 mysql 里面的约束
	UserName   string `json:"user_name"` // 是谁创建的备忘录
	Tittle     string `json:"tittle"`
	Status     int    `json:"status"`      // 0表示List未完成，1为已完成
	Context    string `json:"context"`     // List内容
	StartTime  int64  `json:"start_time"`  // List开始时间
	FinishTime int64  `json:"finish_time"` // List完成时间
}

func BuildTask(task model.Task, uname string) *TaskResponse {
	return &TaskResponse{
		UserName:   uname,
		Tittle:     task.Tittle,
		Status:     task.Status,
		Context:    task.Context,
		StartTime:  task.StartTime,
		FinishTime: task.FinishTime,
	}
}

func BuildTasks(items []model.Task, uname string) *[]TaskResponse {
	var tasks []TaskResponse
	for _, item := range items {
		task := BuildTask(item, uname)
		tasks = append(tasks, *task)
	}
	return &tasks
}
