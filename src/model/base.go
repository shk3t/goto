package model

type SolutionBase struct {
	TaskId int               `json:"taskId"`
	Files  map[string]string `json:"files"`
}

type TaskBase struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	RunTarget   string            `json:"runtarget"`
	Files       map[string]string `json:"files"`
}

type ProjectBase struct {
	Name             string   `json:"name"`
	Language         string   `json:"language"`
	Modules          []string `json:"modules"`
	Containerization string   `json:"containerization"`
	SrcDir           string   `json:"srcdir"`
	StubDir          string   `json:"stubdir"`
}