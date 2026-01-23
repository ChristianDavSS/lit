package config

import (
	"CLI_App/cmd/domain"
	_ "CLI_App/cmd/domain"
	encoder "encoding/json"
	"os"
	"path/filepath"
)

/*
 * json.go - > manages the json files.
 */

// ConfigDto gets the value of the current index on the config file (JSON, etc.).
type ConfigDto struct {
	NamingConventionIndex int8 `json:"activeNamingConventionIndex"`
}

type JsonAdapter struct {
	configPath string
}

func NewJSONAdapter() *JsonAdapter {
	return &JsonAdapter{}
}

func (json *JsonAdapter) GetConfig() *domain.Config {
	var dto *ConfigDto
	err := encoder.Unmarshal(json.readJsonData(), &dto)
	if err != nil {
		os.Exit(1)
	}
	return &domain.Config{
		NamingConventionIndex: dto.NamingConventionIndex,
	}
}

func (json *JsonAdapter) SaveConfig(cfg *ConfigDto) {
	data, err := encoder.Marshal(&cfg)
	if err != nil {
		os.Exit(1)
	}

	err = os.WriteFile(json.GetConfigPath(), data, 0644)
	if err != nil {
		os.Exit(1)
	}
}

func (json *JsonAdapter) GetConfigPath() string {
	if json.configPath != "" {
		return json.configPath
	}

	path, err := os.Executable()
	if err != nil {
		os.Exit(1)
	}
	json.configPath = filepath.Dir(path) + "/config.json"
	return json.configPath
}

func (json *JsonAdapter) readJsonData() []byte {
	file, err := os.ReadFile(json.GetConfigPath())
	if err != nil {
		os.Exit(1)
	}

	return file
}
