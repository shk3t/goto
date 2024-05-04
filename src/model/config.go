package model

import (
	"errors"
	"goto/src/utils"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

func (cfg *GotoConfig) UnmarshalTOML(data any) (fatalError error) {
	fatalError = errors.New("Bad config file format")
	defer func() { recover() }()

	d, _ := data.(map[string]any)

	cfg.Name = d["name"].(string)
	cfg.Language = d["language"].(string)
	cfg.Containerization = utils.GetAssertDefault(d, "containerization", "docker")
	cfg.SrcDir = utils.GetAssertDefault(d, "srcdir", "src")
	cfg.StubDir = utils.GetAssertDefault(d, "stubdir", "stubs")

	packs := d["modules"].([]any)
	for _, p := range packs {
		cfg.Modules = append(cfg.Modules, p.(string))
	}

	taskConfigs := d["tasks"].([]map[string]any)
	cfg.TaskConfigs = make([]TaskConfig, len(taskConfigs))
	taskNames := make([]string, len(taskConfigs))

	for i, tc := range taskConfigs {
		taskConfig := TaskConfig{
			Name:        tc["name"].(string),
			Description: utils.GetAssertDefault(tc, "description", ""),
			RunTarget:   utils.GetAssertDefault(tc, "runtarget", ""),
		}
		cfg.TaskConfigs[i] = taskConfig
		taskNames[i] = tc["name"].(string)

		taskFiles := tc["files"].(any)
		cfg.TaskConfigs[i].Files = map[string]string{}

		switch taskFiles.(type) {
		case []any:
			taskFileNames := make([]string, len(taskFiles.([]any)))

			for j, tf := range taskFiles.([]any) {
				path := tf.(string)
				pathParts := strings.Split(path, string(os.PathSeparator))
				name := pathParts[len(pathParts)-1]
				cfg.TaskConfigs[i].Files[name] = path

				taskFileNames[j] = name
			}

			if !utils.UniqueOnly(&taskFileNames) {
				return errors.New("Conflicting task file names. You should specify them.")
			}
		case map[string]any:
			for name, path := range taskFiles.(map[string]any) {
				cfg.TaskConfigs[i].Files[name] = path.(string)
			}
		}

	}

	if !utils.UniqueOnly(&taskNames) {
		return errors.New("Task names must be unique")
	}

	return nil
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