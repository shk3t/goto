package model

func NewProjectFromConfig(c *GotoConfig) *Project {
	p := Project{}
	p.Name = c.Name
	p.Language = c.Language
	p.Modules = c.Modules
	p.Containerization = c.Containerization
	p.SrcDir = c.SrcDir
	p.StubDir = c.StubDir
	p.Tasks = c.TaskConfigs
	return &p
}