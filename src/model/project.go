package model

type ProjectBase struct {
	Name             string   `json:"name"`
	Language         string   `json:"language"`
	Modules          []string `json:"modules"`
	Containerization string   `json:"containerization"`
	SrcDir           string   `json:"srcdir"`
	StubDir          string   `json:"stubdir"`
}
type Project struct {
	ProjectBase
	Id    int  `json:"id"`
	User  User `json:"user"`
	Dir   string
	Tasks []Task `json:"tasks"`
}
type ProjectPublic struct {
	ProjectBase
	Id    int    `json:"id"`
	Tasks []Task `json:"tasks"`
}
type ProjectMin struct {
	ProjectBase
	Id    int       `json:"id"`
	Tasks []TaskMin `json:"tasks"`
}

func (p *Project) Public() *ProjectPublic {
	return &ProjectPublic{
		ProjectBase: p.ProjectBase,
		Id:          p.Id,
		Tasks:       p.Tasks,
	}
}
func (p *Project) Min() *ProjectMin {
	tasks := make([]TaskMin, len(p.Tasks))
	for i, t := range p.Tasks {
		tasks[i] = *t.Min()
	}

	return &ProjectMin{
		ProjectBase: p.ProjectBase,
		Id:          p.Id,
		Tasks:       tasks,
	}
}

type Module struct {
	Id        int
	ProjectId int
	TaskId    int
	Name      string
}