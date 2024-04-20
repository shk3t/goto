package model

type TaskBase struct {
	Name        string
	Description string
	RunTarget   string
	InjectFiles map[string]string
}
type Task = TaskBase
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
	Url   string
	Dir   string
	Tasks []Task
}
type GotoConfig struct {
	ProjectBase
	TaskConfigs []TaskConfig
}