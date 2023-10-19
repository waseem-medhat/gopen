package main

import (
	"bufio"
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
var configPath = configDir + "/gopenconf"
var config Config

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gopen - a simple CLI to quick-start coding projects")
		return
	}

	switch os.Args[1] {
	case "--get-cmd", "-g":
		checkConfig()
		cmd := getCmd()
		fmt.Println(cmd)

	case "--set-cmd", "-s":
		checkConfig()
		if len(os.Args) < 3 {
			fmt.Println("No command provided")
			return
		}
		setCmd(os.Args[2])

	case "c":
		f, err := os.ReadFile(os.Getenv("HOME") + "/.config/gopen/gopen.json")
		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(f, &config)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(config)

	default:
		gopen(os.Args[1])
	}
}

func checkConfig() {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Println("Config file not found - creating one... ")

		err := os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		_, err = os.Create(configPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func getCmd() string {
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		return scanner.Text()
	}

	return ""
}

func setCmd(cmd string) {
	f, err := os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.WriteString(cmd)
	fmt.Printf("Changed command: %v\n", cmd)
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

	editorCmd := getCmd()
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
