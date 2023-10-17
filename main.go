package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	// "os/exec"
)

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopenconf"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gopen - a simple CLI to quick-start coding projects")
		return
	}

	switch os.Args[1] {
	case "--get-cmd", "-g":
		checkConfig()
        getCmd()

	case "--set-cmd", "-s":
		checkConfig()
        if len(os.Args) < 3 {
            fmt.Println("No command provided")
            return
        }
        setCmd(os.Args[2])
	}

	// "cd /home/waseem/.config/nvim && vi"
	// os.Chdir(os.Getenv("HOME") + "/.config/nvim")
	// cmd := exec.Command(os.Getenv("HOME") + "/Downloads/software/nvim.appimage") // or absolute binary path
	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	// if err != nil {
	// 	panic(err)
	// }
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

func getCmd() {
	f, err := os.OpenFile(configPath, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
    defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		println(scanner.Text())
	}
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
