package model

import "time"

type TaskBase struct {
	TaskConfigBase
	Id        int       `json:"id"`
	ProjectId int       `json:"projectId"`
	Language  string    `json:"language"`
	Modules   []string  `json:"modules"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Task struct {
	TaskBase
	RunTarget string    `json:"runtarget"`
	Files     TaskFiles `json:"files"`
}
type TaskPrivate struct {
	TaskBase
	Files TaskFilesPrivate `json:"files"`
}
type TaskMin struct {
	TaskBase
	FileNames []string `json:"fileNames"`
}

func (t *Task) Private() *TaskPrivate {
	return &TaskPrivate{
		TaskBase: t.TaskBase,
		Files:    t.Files.Private(),
	}
}
func (t *Task) Min() *TaskMin {
	return &TaskMin{
		TaskBase:  t.TaskBase,
		FileNames: t.Files.Names(),
	}
}

type Tasks []Task
type TasksMin []TaskMin

func (tasks Tasks) Min() TasksMin {
	tasksMin := make(TasksMin, len(tasks))
	for i, t := range tasks {
		tasksMin[i] = *t.Min()
	}
	return tasksMin
}

type TaskFileBase struct {
	Id     int    `json:"id"`
	TaskId int    `json:"taskId"`
	Name   string `json:"name"`
	Stub   string `json:"stub"`
}
type TaskFile struct {
	TaskFileBase
	Path string `json:"path"`
}
type TaskFilePrivate struct {
	TaskFileBase
}

func (tf *TaskFile) Private() *TaskFilePrivate {
	return &TaskFilePrivate{TaskFileBase: tf.TaskFileBase}
}

type TaskFiles []TaskFile
type TaskFilesPrivate []TaskFilePrivate

func (taskFiles TaskFiles) Private() TaskFilesPrivate {
	taskFilesPrivate := make(TaskFilesPrivate, len(taskFiles))
	for i, tf := range taskFiles {
		taskFilesPrivate[i] = *tf.Private()
	}
	return taskFilesPrivate
}
func (taskFiles TaskFiles) Names() []string {
	fileNames := make([]string, len(taskFiles))
	for i, tf := range taskFiles {
		fileNames[i] = tf.Name
	}
	return fileNames
}