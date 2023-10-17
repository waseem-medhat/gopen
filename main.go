package main

import (
	"fmt"
	"log"
	"os"
	// "os/exec"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gopen - a simple CLI to quick-start coding projects")
		return
	}

	switch os.Args[1] {
	case "config":
		checkConfig()
	}
	// cmd := exec.Command("vi") // or absolute binary path
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err = cmd.Run()
	// if err != nil {
	//     panic(err)
	// }
}

func checkConfig() {
	configDir := os.Getenv("HOME") + "/.config/gopen"
	configPath := configDir + "/gopen.json"

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
		fmt.Println("Created an empty config file!")
	}

	fmt.Println("Config file found!")
}
