package model

func NewProjectFromConfig(c *GotoConfig) *Project {
	p := Project{ProjectBase: c.ProjectBase}
	p.Tasks = make([]Task, len(c.TaskConfigs))
	for i, tc := range c.TaskConfigs {
		p.Tasks[i] = Task{TaskBase: tc}
	}
	return &p
}