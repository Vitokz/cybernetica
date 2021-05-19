package model

import "time"

type Task struct {
	tableName    interface{} `sql:"tasks"`
	Task_id      string      `json:"task_id"`
	Group_id     int         `json:"group_id"`
	Task         string      `json:"task"`
	Completed    bool        `json:"completed" pg:",use_zero"`
	Created_at   time.Time   `json:"created_at"`
	Completed_at time.Time   `json:"completed_at"`
}

type TaskNew struct {
	tableName interface{} `sql:"tasks"`
	Group_id  int         `json:"group_id"`
	Task      string      `json:"task"`
}

type TaskStat struct {
	Completed int `json:"completed"`
	Created   int `json:"created"`
}
