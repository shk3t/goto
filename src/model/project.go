package model

import "time"

type ProjectBase struct {
	ProjectConfigBase
	Id        int       `json:"id"`
	UserId    int       `json:"userId"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Project struct {
	ProjectBase
	Dir   string
	Tasks Tasks `json:"tasks"`
}
type ProjectPublic struct {
	ProjectBase
	Tasks Tasks `json:"tasks"`
}
type ProjectMin struct {
	ProjectBase
	Tasks TasksMin `json:"tasks"`
}

func (p *Project) Public() *ProjectPublic {
	return &ProjectPublic{
		ProjectBase: p.ProjectBase,
		Tasks:       p.Tasks,
	}
}
func (p *Project) Min() *ProjectMin {
	tasks := make(TasksMin, len(p.Tasks))
	for i, t := range p.Tasks {
		tasks[i] = *t.Min()
	}

	return &ProjectMin{
		ProjectBase: p.ProjectBase,
		Tasks:       tasks,
	}
}

type Projects []Project
type ProjectsMin []ProjectMin

func (projects Projects) Min() ProjectsMin {
	projectsMin := make(ProjectsMin, len(projects))
	for i, p := range projects {
		projectsMin[i] = *p.Min()
	}
	return projectsMin
}

type Module struct {
	Id        int
	ProjectId int
	TaskId    int
	Name      string
}

type Modules []Module

func (modules Modules) Names() []string {
	names := make([]string, len(modules))
	for i, m := range modules {
		names[i] = m.Name
	}
	return names
}