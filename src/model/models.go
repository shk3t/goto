package model

import "time"

type Solution struct {
	Id        int       `json:"id"`
	UserId    int       `json:"userId"`
	TaskId    int       `json:"taskId"`
	Status    string    `json:"status"`
	Code      string    `json:"code"`
	Result    string    `json:"result"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type User struct {
	Id       int
	Login    string
	Password string
	IsAdmin  bool
}

type Task struct {
	TaskBase
	Id        int `json:"id"`
	ProjectId int `json:"projectId"`
}
type TaskConfig = TaskBase

type Project struct {
	ProjectBase
	Id    int
	User  User
	Dir   string
	Tasks []Task
}
type GotoConfig struct {
	ProjectBase
	TaskConfigs []TaskConfig
}

type InjectFile struct {
	Id     int
	TaskId int
	Name   string
	Path   string
}