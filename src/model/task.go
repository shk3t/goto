package model

type TaskBase struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Modules     []string `json:"modules"`
}
type Task struct {
	TaskBase
	ProjectId int        `json:"projectId"`
	RunTarget string     `json:"runtarget"`
	Files     []TaskFile `json:"files"`
}
type TaskPrivate struct {
	TaskBase
	Description string            `json:"description"`
	Files       []TaskFilePrivate `json:"files"`
}
type TaskMin struct {
	TaskBase
	ProjectId int      `json:"projectId"`
	FileNames []string `json:"fileNames"`
}

func (t *Task) Private() *TaskPrivate {
	taskFiles := make([]TaskFilePrivate, len(t.Files))
	for i, tf := range t.Files {
		taskFiles[i] = *tf.Private()
	}

	return &TaskPrivate{
		TaskBase:    t.TaskBase,
		Description: t.Description,
		Files:       taskFiles,
	}
}
func (t *Task) Min() *TaskMin {
	fileNames := make([]string, len(t.Files))
	for i, tf := range t.Files {
		fileNames[i] = tf.Name
	}

	return &TaskMin{
		TaskBase:  t.TaskBase,
		ProjectId: t.ProjectId,
		FileNames: fileNames,
	}
}

type TaskFile struct {
	Id     int    `json:"id"`
	TaskId int    `json:"taskId"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Stub   string `json:"stub"`
}
type TaskFilePrivate struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Stub string `json:"stub"`
}

func (tf *TaskFile) Private() *TaskFilePrivate {
	return &TaskFilePrivate{
		Id:   tf.Id,
		Name: tf.Name,
		Stub: tf.Stub,
	}
}