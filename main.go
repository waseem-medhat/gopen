package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	EditorCmd string  `json:"editorCmd"`
	Aliases   []Alias `json:"aliases"`
}

type Alias struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopen.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gopen - a simple CLI to quick-start coding projects")
		return
	}

	switch os.Args[1] {
	case "--init", "-i":
		initConfig(configDir, configPath)

	case "--get-cmd", "-g":
		config := readConfig(configPath)
		fmt.Println(config.EditorCmd)

	case "--set-cmd", "-s":
		if len(os.Args) < 3 {
			fmt.Println("No command provided")
			return
		}
		setCmd(os.Args[2])

	default:
		gopen(os.Args[1])
	}
}

func initConfig(configDir string, configPath string) {
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

	emptyConfig := Config{"", []Alias{}}
	writeConfig(emptyConfig, configPath)
}

func writeConfig(config Config, configPath string) {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(configPath, jsonFile, 0644)
}

func readConfig(configPath string) Config {
	var config Config

	f, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func setCmd(cmd string) {
	config := readConfig(configPath)
	config.EditorCmd = cmd
	writeConfig(config, configPath)
}

func gopen(path string) {
	fInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		fmt.Println("Path doesn't exist")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	if !fInfo.IsDir() {
		println("Not a directory")
		return
	}

	config := readConfig(configPath)
	editorCmd := config.EditorCmd
	os.Chdir(path)
	cmd := exec.Command(editorCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
