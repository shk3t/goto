package model

import "time"

type Solution struct {
	SolutionBase
	Id        int               `json:"id"`
	UserId    int               `json:"userId"`
	Status    string            `json:"status"`
	Result    string            `json:"result"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

type SolutionFile struct {
	Id         int
	SolutionId int
	Name       string
	Code       string
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

type TaskFile struct {
	Id     int
	TaskId int
	Name   string
	Path   string
}