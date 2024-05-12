package model

import (
	"errors"
	u "goto/src/utils"
	"os"
	sc "strconv"
	s "strings"

	"github.com/BurntSushi/toml"
)

type TaskConfigBase struct {
	Name        string
	Description string
}
type TaskConfig struct {
	TaskConfigBase
	RunTarget string
	Files     TaskFiles
	OldName   string
}

func (tc *TaskConfig) Task() *Task {
	return &Task{
		TaskBase:  TaskBase{TaskConfigBase: tc.TaskConfigBase},
		RunTarget: tc.RunTarget,
		Files:     tc.Files,
	}
}

type TaskConfigs []TaskConfig

type ProjectConfigBase struct {
	Name             string
	Language         string
	Modules          []string
	Containerization string
	SrcDir           string
	StubDir          string
}
type GotoConfig struct {
	ProjectConfigBase
	TaskConfigs TaskConfigs
}

func (cfg *GotoConfig) Project() *Project {
	p := Project{}
	p.ProjectConfigBase = cfg.ProjectConfigBase
	p.Tasks = make(Tasks, len(cfg.TaskConfigs))
	for i, tc := range cfg.TaskConfigs {
		p.Tasks[i] = *tc.Task()
	}
	return &p
}

func LoadGotoConfig(configPath string) (*GotoConfig, error) {
	var config GotoConfig

	tomlBytes, err := os.ReadFile(configPath)
	if err != nil {
		return &config, err
	}

	_, err = toml.Decode(string(tomlBytes), &config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}

func (cfg *GotoConfig) UnmarshalTOML(data any) (fatalError error) {
	var err error
	var ok bool
	fatalError = errors.New("Bad config file format")
	defer func() { recover() }()

	d, _ := data.(map[string]any)

	cfg.Name, err = u.GetAssertError[string](d, "name", "")
	if err != nil {
		return err
	}
	cfg.Language, err = u.GetAssertError[string](d, "language", "")
	if err != nil {
		return err
	}
	cfg.Containerization = u.GetAssertDefault(d, "containerization", "docker")
	cfg.SrcDir = u.GetAssertDefault(d, "srcdir", "src")
	cfg.StubDir = u.GetAssertDefault(d, "stubdir", "stubs")

	modules, err := u.GetAssertError[[]any](d, "modules", "")
	if err != nil {
		return err
	}

	cfg.Modules = make([]string, len(modules))
	for i, m := range modules {
		cfg.Modules[i], ok = m.(string)
		if !ok {
			return errors.New("`modules` has bad format")
		}
	}

	taskConfigs, err := u.GetAssertError[[]map[string]any](d, "tasks", "")
	if err != nil {
		return err
	}
	cfg.TaskConfigs = make(TaskConfigs, len(taskConfigs))
	taskNames := make([]string, len(taskConfigs))

	for i, tc := range taskConfigs {
		taskName, err := u.GetAssertError[string](tc, "name", sc.Itoa(i+1)+" task")
		if err != nil {
			return err
		}
		taskConfig := TaskConfig{
			TaskConfigBase: TaskConfigBase{
				Name:        taskName,
				Description: u.GetAssertDefault(tc, "description", ""),
			},
			RunTarget: u.GetAssertDefault(tc, "runtarget", ""),
			OldName:   u.GetAssertDefault(tc, "oldname", ""),
		}
		cfg.TaskConfigs[i] = taskConfig
		taskNames[i] = taskName

		taskFiles, err := u.GetAssertError[any](tc, "files", taskName+" task")
		cfg.TaskConfigs[i].Files = TaskFiles{}

		switch taskFiles.(type) {
		case []any:
			taskFiles := taskFiles.([]any)
			taskFileNames := make([]string, len(taskFiles))
			cfg.TaskConfigs[i].Files = make(TaskFiles, len(taskFiles))

			for j, tf := range taskFiles {
				path, ok := tf.(string)
				if !ok {
					return errors.New(taskName + " task, " + sc.Itoa(j+1) + " file: `path` has bad format")
				}
				pathParts := s.Split(path, string(os.PathSeparator))
				name := pathParts[len(pathParts)-1]

				task := TaskFile{TaskFileBase: TaskFileBase{Name: name}, Path: path}
				cfg.TaskConfigs[i].Files[j] = task

				taskFileNames[j] = name
			}

			if !u.UniqueOnly(&taskFileNames) {
				return errors.New(taskName + " task: conflicting file names; you should specify them via hashtable.")
			}
		case map[string]any:
			taskFiles := taskFiles.(map[string]any)
			for name, path := range taskFiles {
				path, ok := path.(string)
				if !ok {
					return errors.New(taskName + " task, " + name + " file: `path` has bad format")
				}
				task := TaskFile{TaskFileBase: TaskFileBase{Name: name}, Path: path}
				cfg.TaskConfigs[i].Files = append(cfg.TaskConfigs[i].Files, task)
			}
		default:
			return errors.New(taskName + " task: `files` has bad format")
		}

	}

	if !u.UniqueOnly(&taskNames) {
		return errors.New("Task names must be unique")
	}

	return nil
}