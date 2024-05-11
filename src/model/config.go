package model

import (
	"errors"
	"goto/src/utils"
	"os"
	sc "strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type TaskConfigBase struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type TaskConfig struct {
	TaskConfigBase
	RunTarget string     `json:"runtarget"`
	Files     []TaskFile `json:"files"`
}

func (tc *TaskConfig) Task() *Task {
	return &Task{
		TaskBase:  TaskBase{TaskConfigBase: tc.TaskConfigBase},
		RunTarget: tc.RunTarget,
		Files:     tc.Files,
	}
}

type ProjectConfigBase struct {
	Name             string   `json:"name"`
	Language         string   `json:"language"`
	Modules          []string `json:"modules"`
	Containerization string   `json:"containerization"`
	SrcDir           string   `json:"srcdir"`
	StubDir          string   `json:"stubdir"`
}
type GotoConfig struct {
	ProjectConfigBase
	TaskConfigs []TaskConfig
}

func (cfg *GotoConfig) Project() *Project {
	p := Project{ProjectBase: ProjectBase{ProjectConfigBase: cfg.ProjectConfigBase}}
	p.Tasks = make([]Task, len(cfg.TaskConfigs))
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

	cfg.Name, err = utils.GetAssertError[string](d, "name", "")
	if err != nil {
		return err
	}
	cfg.Language, err = utils.GetAssertError[string](d, "language", "")
	if err != nil {
		return err
	}
	cfg.Containerization = utils.GetAssertDefault(d, "containerization", "docker")
	cfg.SrcDir = utils.GetAssertDefault(d, "srcdir", "src")
	cfg.StubDir = utils.GetAssertDefault(d, "stubdir", "stubs")

	modules, err := utils.GetAssertError[[]any](d, "modules", "")
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

	taskConfigs, err := utils.GetAssertError[[]map[string]any](d, "tasks", "")
	if err != nil {
		return err
	}
	cfg.TaskConfigs = make([]TaskConfig, len(taskConfigs))
	taskNames := make([]string, len(taskConfigs))

	for i, tc := range taskConfigs {
		taskName, err := utils.GetAssertError[string](tc, "name", sc.Itoa(i+1)+" task")
		if err != nil {
			return err
		}
		taskConfig := TaskConfig{
			TaskConfigBase: TaskConfigBase{
				Name:        taskName,
				Description: utils.GetAssertDefault(tc, "description", ""),
			},
			RunTarget: utils.GetAssertDefault(tc, "runtarget", ""),
		}
		cfg.TaskConfigs[i] = taskConfig
		taskNames[i] = taskName

		taskFiles, err := utils.GetAssertError[any](tc, "files", taskName+" task")
		cfg.TaskConfigs[i].Files = []TaskFile{}

		switch taskFiles.(type) {
		case []any:
			taskFiles := taskFiles.([]any)
			taskFileNames := make([]string, len(taskFiles))
			cfg.TaskConfigs[i].Files = make([]TaskFile, len(taskFiles))

			for j, tf := range taskFiles {
				path, ok := tf.(string)
				if !ok {
					return errors.New(taskName + " task, " + sc.Itoa(j+1) + " file: `path` has bad format")
				}
				pathParts := strings.Split(path, string(os.PathSeparator))
				name := pathParts[len(pathParts)-1]

				task := TaskFile{Name: name, Path: path}
				cfg.TaskConfigs[i].Files[j] = task

				taskFileNames[j] = name
			}

			if !utils.UniqueOnly(&taskFileNames) {
				return errors.New(taskName + " task: conflicting file names; you should specify them via hashtable.")
			}
		case map[string]any:
			taskFiles := taskFiles.(map[string]any)
			for name, path := range taskFiles {
				path, ok := path.(string)
				if !ok {
					return errors.New(taskName + " task, " + name + " file: `path` has bad format")
				}
				task := TaskFile{Name: name, Path: path}
				cfg.TaskConfigs[i].Files = append(cfg.TaskConfigs[i].Files, task)
			}
		default:
			return errors.New(taskName + " task: `files` has bad format")
		}

	}

	if !utils.UniqueOnly(&taskNames) {
		return errors.New("Task names must be unique")
	}

	return nil
}