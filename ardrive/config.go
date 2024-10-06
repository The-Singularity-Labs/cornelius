package ardrive

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Id             string `yaml:"id"`
	WalletPath     string `yaml:"wallet_path"`
	Password       string `yaml:"password"`
	ParentFolderId string `yaml:"parent_folder_id"`
	IsPublic       bool   `yaml:"is_public"`
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
