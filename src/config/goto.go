package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type GotoConfig struct {
	Name             string
	Language         string
	Packages         []string
	Containerization string
	SrcDir           string
	StubDir          string
}

func LoadGotoConfig(configPath string) (GotoConfig, error) {
	var config GotoConfig

	tomlBytes, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	if _, err := toml.Decode(string(tomlBytes), &config); err != nil {
		return config, err
	}
	return config, nil
}