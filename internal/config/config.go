// Package config includes functions for initializing, reading, and writing
// Gopen config files.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// InitConfig checks if the config file exists in configPath. If not, creates
// an empty config file. configDir will also be created if it doesn't exist.
func InitConfig(configDir string, configPath string) {
	if _, err := os.Stat(configPath); err == nil {
		fmt.Println("Found config file - exiting...")
		return
	}

	fmt.Println("Creating a new config file...")
	err := os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	_, err = os.Create(configPath)
	if err != nil {
		log.Fatal(err)
	}

	emptyConfig := structs.Config{}
	WriteConfig(emptyConfig, configPath)
}

// WriteConfig writes config to configPath (will OVERWRITE if file already
// exists)
func WriteConfig(config structs.Config, configPath string) {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(configPath, jsonFile, 0644)
}

// ReadConfig reads the configPath file and returns a Config struct
func ReadConfig(configPath string) (structs.Config, error) {
	var config structs.Config

	f, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
