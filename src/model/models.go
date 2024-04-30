package model

type User struct {
	Id       int
	Login    string
	Password string
	IsAdmin  bool
}

type TaskBase struct {
	Name        string
	Description string
	RunTarget   string
	InjectFiles map[string]string
}
type Task struct {
	TaskBase
	Id int
}
type TaskConfig = TaskBase

type ProjectBase struct {
	Name             string
	Language         string
	Modules          []string
	Containerization string
	SrcDir           string
	StubDir          string
}
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