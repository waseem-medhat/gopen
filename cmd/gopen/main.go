package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/diralias"
	"github.com/wipdev-tech/gopen/internal/gopen"
	"github.com/wipdev-tech/gopen/internal/structs"
)

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopen.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Gopen - a simple CLI to quick-start coding projects")
		return
	}

	switch os.Args[1] {
	case "--init", "-i":
		config.InitConfig(configDir, configPath)

	case "--editor-cmd", "-e":
		configObj, err := config.ReadConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}

		if len(os.Args) < 3 {
			fmt.Println(configObj.EditorCmd)
		} else {
			configObj.EditorCmd = os.Args[2]
			err := config.WriteConfig(configObj, configPath)
			if err != nil {
				log.Fatal(err)
			}
		}

	case "--alias", "-a":
		configObj, err := config.ReadConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}

		if len(os.Args) < 3 {
			diralias.ListDirAliases(configObj)
		} else if len(os.Args) > 4 {
			fmt.Println("Too many arguments - exiting...")
		} else {
			configObj.DirAliases = append(
				configObj.DirAliases,
				structs.DirAlias{Alias: os.Args[2], Path: os.Args[3]},
			)
			err := config.WriteConfig(configObj, configPath)
			if err != nil {
				log.Fatal(err)
			}
		}

	default:
		configObj, err := config.ReadConfig(configPath)
		if err != nil {
			log.Fatal(err)
		}

		gopen.Gopen(os.Args[1], configObj)
	}
}
