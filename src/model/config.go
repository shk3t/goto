package model

import (
	"errors"
	"fmt"
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
	cfg.Containerization = d["containerization"].(string)
	cfg.SrcDir = d["srcdir"].(string)
	cfg.StubDir = d["stubdir"].(string)

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
			Description: tc["description"].(string),
			RunTarget:   tc["runtarget"].(string),
		}
		cfg.TaskConfigs[i] = taskConfig
		taskNames[i] = tc["name"].(string)

		injectFiles := tc["injectfiles"].(any)
		cfg.TaskConfigs[i].InjectFiles = make(map[string]string)

		switch injectFiles.(type) {
		case []any:
			injectFileNames := make([]string, len(injectFiles.([]any)))

			for j, ifl := range injectFiles.([]any) {
				path := ifl.(string)
				pathParts := strings.Split(path, string(os.PathSeparator))
				name := pathParts[len(pathParts)-1]
				cfg.TaskConfigs[i].InjectFiles[name] = path

				injectFileNames[j] = name
			}

			if !utils.UniqueOnly(&injectFileNames) {
				return errors.New("Conflicting InjectFile names. You should specify them.")
			}
		case map[string]any:
			for name, path := range injectFiles.(map[string]any) {
				cfg.TaskConfigs[i].InjectFiles[name] = path.(string)
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
	fmt.Println(config)
	if err != nil {
		return &config, err
	}
	return &config, nil
}