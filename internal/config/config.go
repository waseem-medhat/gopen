// Package config includes functions for initializing, reading, and writing
// Gopen config files.
package config

import (
	"encoding/json"
	"os"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// Init checks if the config file exists in configPath. If not, creates an
// empty config file. configDir will also be created if it doesn't exist.
func Init(configDir string, configPath string) (err error) {
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
	err = Write(emptyConfig, configPath)
	if err != nil {
		return
	}

	return
}

// Write writes config to configPath (will OVERWRITE if file already exists)
func Write(config structs.Config, configPath string) (err error) {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

func Migrate(configPath string) error {
    originalConfig, err := Read(configPath)
    if err != nil {
        return err 
    }


    var newConfig = &structs.Config{
        EditorCmd: originalConfig.EditorCmd,
        CustomBehaviour: false,
        DirAliases: originalConfig.DirAliases,
    }

    err = Write(*newConfig, configPath)
    if err != nil {
        return err
    }

    return nil
}

// Read reads the configPath file and returns a Config struct
func Read(configPath string) (config structs.Config, err error) {
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
