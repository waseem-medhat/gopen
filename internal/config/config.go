// Package config includes functions for initializing, reading, and writing
// Gopen config files.
package config

import (
	"encoding/json"
	"os"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// InitConfig checks if the config file exists in configPath. If not, creates
// an empty config file. configDir will also be created if it doesn't exist.
func InitConfig(configDir string, configPath string) (err error) {
	_, err = os.Stat(configPath)
	if err == nil {
		return os.ErrExist
	}

	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return
	}

	_, err = os.Create(configPath)
	if err != nil {
		return
	}

	emptyConfig := structs.Config{}
	err = WriteConfig(emptyConfig, configPath)
	if err != nil {
		return
	}

	return
}

// WriteConfig writes config to configPath (will OVERWRITE if file already
// exists)
func WriteConfig(config structs.Config, configPath string) (err error) {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return
	}

	err = os.WriteFile(configPath, jsonFile, 0644)
	if err != nil {
		return
	}

	return
}

// ReadConfig reads the configPath file and returns a Config struct
func ReadConfig(configPath string) (config structs.Config, err error) {
	f, err := os.ReadFile(configPath)
	if err != nil {
		return
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		return
	}

	return
}
