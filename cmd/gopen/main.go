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
		fmt.Println("Gopen - a simple CLI to quick-start coding projects")
		return
	}

	switch os.Args[1] {
	case "--init", "-i":
		err := config.InitConfig(configDir, configPath)
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
		}

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
			for _, fmtAlias := range diralias.ListDirAliases(configObj) {
				fmt.Println(fmtAlias)
			}
		} else if len(os.Args) > 4 {
			fmt.Println("Too many arguments - exiting...")
		} else {
			configObj = diralias.Add(configObj, os.Args[2], os.Args[3])
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
