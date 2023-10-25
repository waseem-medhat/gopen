package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	EditorCmd  string     `json:"editorCmd"`
	DirAliases []DirAlias `json:"aliases"`
}

type DirAlias struct {
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

	case "--editor-cmd", "-e":
		config := readConfig(configPath)

		if len(os.Args) < 3 {
			fmt.Println(config.EditorCmd)
		} else {
			config.EditorCmd = os.Args[2]
			writeConfig(config, configPath)
		}

	case "--alias", "-a":
		config := readConfig(configPath)
		if len(os.Args) < 3 {
			listDirAliases(config)
		} else if len(os.Args) > 4 {
			fmt.Println("Too many arguments - exiting...")
		} else {
			config.DirAliases = append(
				config.DirAliases,
				DirAlias{os.Args[2], os.Args[3]},
			)
			writeConfig(config, configPath)
		}

	default:
		config := readConfig(configPath)
		gopen(os.Args[1], config)
	}
}

// initConfig checks if the config file exists in configPath. If not, creates
// an empty config file. configDir will also be created if it doesn't exist.
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

	emptyConfig := Config{"", []DirAlias{}}
	writeConfig(emptyConfig, configPath)
}

// writeConfig writes config to configPath (will OVERWRITE if file already
// exists)
func writeConfig(config Config, configPath string) {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile(configPath, jsonFile, 0644)
}

// readConfig reads the configPath file and returns a Config struct
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

// listDirAliases pretty-prints each alias and its corresponding path
func listDirAliases(config Config) {
	var width int
	for _, dirAlias := range config.DirAliases {
		if len(dirAlias.Alias) > width {
			width = len(dirAlias.Alias)
		}
	}

	for _, dirAlias := range config.DirAliases {
		fmt.Printf("%*s: %s\n", width, dirAlias.Alias, dirAlias.Path)
	}
}

// gopen uses the Config struct to find the path corresponding to targetAlias
// and executes the editor command with the target path as the working
// directory
func gopen(targetAlias string, config Config) {
	var targetPath string
	for _, dirAlias := range config.DirAliases {
		if targetAlias == dirAlias.Alias {
			targetPath = dirAlias.Path
			break
		}
	}

	fInfo, err := os.Stat(targetPath)
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

	editorCmd := config.EditorCmd
	os.Chdir(targetPath)
	cmd := exec.Command(editorCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
