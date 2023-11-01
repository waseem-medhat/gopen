package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/diralias"
	"github.com/wipdev-tech/gopen/internal/gopen"
)

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopen.json"

func main() {
	if len(os.Args) < 2 {
		handleHelp()
		return
	}

	switch os.Args[1] {
	case "--init", "-i":
		handleInit()

	case "--editor-cmd", "-e":
		handleEditorCmd()

	case "--alias", "-a":
		handleAlias()

	default:
		handleGopen()
	}
}

func handleHelp() {
	fmt.Println("Gopen - a simple CLI to quick-start coding projects")
}

func handleInit() {
	err := config.Init(configDir, configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}
}

func handleEditorCmd() {
	configObj, err := config.Read(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) < 3 {
		fmt.Println(configObj.EditorCmd)
	} else {
		configObj.EditorCmd = os.Args[2]
		err := config.Write(configObj, configPath)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func handleAlias() {
	configObj, err := config.Read(configPath)
	if err != nil {
		log.Fatal(err)
	}

	switch len(os.Args) {
	case 2:
		for _, fmtAlias := range diralias.List(configObj) {
			fmt.Println(fmtAlias)
		}

	case 3:
		for _, dirAlias := range configObj.DirAliases {
			if dirAlias.Alias == os.Args[2] {
				fmt.Println(dirAlias.Path)
			}
		}

	case 4:
		configObj = diralias.Add(configObj, os.Args[2], os.Args[3])
		err := config.Write(configObj, configPath)
		if err != nil {
			log.Fatal(err)
		}

	default:
		fmt.Println("Too many arguments - exiting...")
	}
}

func handleGopen() {
	configObj, err := config.Read(configPath)
	if err != nil {
		log.Fatal(err)
	}

	gopen.Gopen(os.Args[1], configObj)
}
