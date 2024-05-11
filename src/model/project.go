package model

type ProjectBase struct {
	ProjectConfigBase
	Id     int `json:"id"`
	UserId int `json:"userId"`
}

type Project struct {
	ProjectBase
	Dir   string
	Tasks []Task `json:"tasks"`
}
type ProjectPublic struct {
	ProjectBase
	Tasks []Task `json:"tasks"`
}
type ProjectMin struct {
	ProjectBase
	Tasks []TaskMin `json:"tasks"`
}

func (p *Project) Public() *ProjectPublic {
	return &ProjectPublic{
		ProjectBase: p.ProjectBase,
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
		Tasks:       tasks,
	}
}

type Module struct {
	Id        int
	ProjectId int
	TaskId    int
	Name      string
}