package sync

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Concurrency  int        `yaml:"concurrency"`
	TmpDirectory string     `yaml:"tmp_directory"`
	Pipelines    []Pipeline `yaml:"pipelines"`
}

func LoadConfig(path string) (Config, error) {
	cfg := Config{}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("unable to read file path %q: %w", path, err)
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		fmt.Errorf("unable to parse config yaml: %w", err)
	}

	return cfg, nil
}
