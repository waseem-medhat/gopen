package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/diralias"
	"github.com/wipdev-tech/gopen/internal/gopen"
	"github.com/wipdev-tech/gopen/internal/structs"
)

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopen.json"

func main() {
	if len(os.Args) < 2 {
		handleHelp()
		return
	}

	switch os.Args[1] {
	case "help", "h":
		handleHelp()

	case "init", "i":
		handleInit()

	case "editor", "e":
		handleEditor()

	case "alias", "a":
		handleAlias()

	case "remove", "r":
		handleRemove()
	case "migrate", "m":
		handleMigrate()
	case "custom", "c":
		handleCustom()
	default:
		handleGopen()
	}
}

func handleInit() {
	err := config.Init(configDir, configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}
}

func handleEditor() {
	configObj, err := config.Read(configPath)
	errFatal(err)

	if len(os.Args) < 3 {
		fmt.Println(configObj.EditorCmd)
	} else {
		configObj.EditorCmd = os.Args[2]
		err := config.Write(configObj, configPath)
		errFatal(err)
	}
}

func handleAlias() {
	configObj, err := config.Read(configPath)
	errFatal(err)

	switch len(os.Args) {
	case 2:
		for _, fmtAlias := range diralias.List(configObj) {
			fmt.Println(fmtAlias)
		}

	case 3:
		for _, dirAlias := range configObj.DirAliases {
			if dirAlias.Alias == os.Args[2] {
				fmt.Println(dirAlias.Path)
				return
			}
		}
		fmt.Println("Alias doesn't exist")

	case 4:
		configObj, err := diralias.Add(configObj, os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		err = config.Write(configObj, configPath)
		errFatal(err)

	default:
		fmt.Println("Too many arguments - exiting...")
	}
}

func handleGopen() {
	configObj, err := config.Read(configPath)
	errFatal(err)

	err = gopen.Gopen(os.Args[1], configObj)
	if err != nil {
		fmt.Println(err)
	}
}

func errFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleRemove() {
	configObj, err := config.Read(configPath)
	errFatal(err)

	if len(os.Args) != 3 {
		fmt.Println("Error: must provide one alias to 'remove' command")
		return
	}

	var newConfig structs.Config
	newConfig.EditorCmd = configObj.EditorCmd
	for _, dirAlias := range configObj.DirAliases {
		if dirAlias.Alias != os.Args[2] {
			newConfig.DirAliases = append(newConfig.DirAliases, dirAlias)
		}
	}

	err = config.Write(newConfig, configPath)
	errFatal(err)
}

func handleMigrate() {
	err := config.Migrate(configPath)
	if err != nil {
		errFatal(err)
	}
	fmt.Println("Successfully migrated config")
}
func handleCustom() {

	if len(os.Args) < 3 {
		fmt.Println("Unexpected number of args: expected 2")
	}

	arg := strings.ToLower(os.Args[2])

	if arg != "false" && arg != "true" {
		err := fmt.Errorf("Error: expected argument true or false, got %v", arg)
		errFatal(err)
	}

	custom := true

	if arg == "false" {
		custom = false
	}
	configObj, err := config.Read(configPath)
	if err != nil {
		errFatal(err)
	}

	configObj.CustomBehaviour = custom
	err = config.Write(configObj, configPath)
	if err != nil {
		errFatal(err)
	}

	fmt.Printf("Successfully set custom behaviour to: %s\n", arg)

}
func handleHelp() {
	fmt.Print(`Gopen - a simple CLI to quick-start coding projects

Usage:

    gopen foo         cd into path assigned to alias 'foo' and run the editor cmd
    gopen cmd [args]  Run command 'cmd' (see Commands below)

Commands:
Can be abbreviated by the first letter ('gopen i' == 'gopen init')

    init              Initialize a new config file (~/.config/gopen/gopen.json)

    editor            Get editor command
    editor cmd        Set editor command to 'cmd'

    alias             List all saved aliases
    alias foo         Get path assigned to alias 'foo'
    alias foo bar     Assign to alias 'foo' the path 'bar'

    remove foo        Remove alias 'foo' from the config

    help              Print this help message

`)
}
