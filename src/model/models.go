package model

type ProjectBase struct {
	Name             string
	Language         string
	Modules          []string
	Containerization string
	SrcDir           string
	StubDir          string
}

type TaskBase struct {
	Name        string
	Description string
	RunTarget   string
	InjectFiles map[string]string
}

type Project struct {
	ProjectBase
	Id    int
	Url   string
	Dir   string
	Tasks []Task
}

type Task = TaskBase