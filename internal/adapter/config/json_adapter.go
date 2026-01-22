package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"CLI_App/internal/domain"
)

// ConfigDTO is the data transfer object for the configuration JSON.
type ConfigDTO struct {
	NamingConventionIndex int8 `json:"activeNamingConventionIndex"`
}

// JSONConfigAdapter implements the ConfigProvider interface using a JSON file.
type JSONConfigAdapter struct {
	configPath string
}

// NewJSONConfigAdapter creates a new JSONConfigAdapter.
func NewJSONConfigAdapter() *JSONConfigAdapter {
	return &JSONConfigAdapter{}
}

// getConfigPath returns the path to the config file.
func (a *JSONConfigAdapter) getConfigPath() (string, error) {
	if a.configPath != "" {
		return a.configPath, nil
	}
	path, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("error getting executable path: %w", err)
	}
	a.configPath = filepath.Dir(path) + "/config.json"
	return a.configPath, nil
}

// GetConfig reads the configuration from the JSON file.
func (a *JSONConfigAdapter) GetConfig() (*domain.Config, error) {
	path, err := a.getConfigPath()
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(path)
	if err != nil {
		// If file doesn't exist, return default config?
		// For now, return error as per original behavior, but maybe wrapped.
		return nil, fmt.Errorf("couldn't read config file at %s: %w", path, err)
	}

	var dto ConfigDTO
	if err := json.Unmarshal(file, &dto); err != nil {
		return nil, fmt.Errorf("error parsing config json: %w", err)
	}

	return &domain.Config{
		NamingConventionIndex: dto.NamingConventionIndex,
	}, nil
}

// SaveConfig writes the configuration to the JSON file.
func (a *JSONConfigAdapter) SaveConfig(cfg *domain.Config) error {
	path, err := a.getConfigPath()
	if err != nil {
		return err
	}

	dto := ConfigDTO{
		NamingConventionIndex: cfg.NamingConventionIndex,
	}

	data, err := json.Marshal(dto)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}
