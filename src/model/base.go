package model

type TaskBase struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	RunTarget   string            `json:"runtarget"`
	InjectFiles map[string]string `json:"injectfiles"`
}

type ProjectBase struct {
	Name             string   `json:"name"`
	Language         string   `json:"language"`
	Modules          []string `json:"modules"`
	Containerization string   `json:"containerization"`
	SrcDir           string   `json:"srcdir"`
	StubDir          string   `json:"stubdir"`
}