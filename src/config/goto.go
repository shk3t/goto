package config

import (
	"errors"
	"fmt"
	"goto/src/model"
	"os"
	"strings"

	"github.com/BurntSushi/toml"
)

type GotoConfig struct {
	model.ProjectBase
	TaskConfigs []TaskConfig
}

type TaskConfig = model.TaskBase

func (cfg *GotoConfig) UnmarshalTOML(data any) (fatalError error) {
	fatalError = errors.New("Bad config file format")
	defer func() { recover() }()

	d, _ := data.(map[string]any)

	cfg.Name = d["name"].(string)
	cfg.Language = d["language"].(string)
	cfg.Containerization = d["containerization"].(string)
	cfg.SrcDir = d["srcdir"].(string)
	cfg.StubDir = d["stubdir"].(string)

	packs := d["packages"].([]any)
	for _, p := range packs {
		cfg.Modules = append(cfg.Modules, p.(string))
	}

	cfg.TaskConfigs = make([]TaskConfig, 3)
	taskConfigs := d["tasks"].([]map[string]any)

	for i, tc := range taskConfigs {
		taskConfig := TaskConfig{
			Name:        tc["name"].(string),
			Description: tc["description"].(string),
			RunTarget:   tc["runtarget"].(string),
		}
		cfg.TaskConfigs[i] = taskConfig

		cfg.TaskConfigs[i].InjectFiles = make(map[string]string)
		injectFiles := tc["injectfiles"].(any)

		switch injectFiles.(type) {
		case []any:
			for _, ifl := range injectFiles.([]any) {
				filePath := ifl.(string)
				filePathParts := strings.Split(filePath, string(os.PathSeparator))
				fileName := filePathParts[len(filePathParts)-1]
				cfg.TaskConfigs[i].InjectFiles[fileName] = filePath
			}
		case map[string]any:
			for fileName, filePath := range injectFiles.(map[string]any) {
				cfg.TaskConfigs[i].InjectFiles[fileName] = filePath.(string)
			}
		}
	}

	return nil
}

func LoadGotoConfig(configPath string) (GotoConfig, error) {
	var config GotoConfig

	tomlBytes, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	_, err = toml.Decode(string(tomlBytes), &config)
	fmt.Println(config)
	if err != nil {
		return config, err
	}
	return config, nil
}