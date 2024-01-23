package main

import (
	"fmt"
	"log"
	"os"

	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/diralias"
	"github.com/wipdev-tech/gopen/internal/fzf"
	"github.com/wipdev-tech/gopen/internal/gopen"
	"github.com/wipdev-tech/gopen/internal/structs"
)

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopen.json"

func main() {
	if len(os.Args) < 2 {
		handleFuzzy()
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
func handleCustom() {
	configObj, err := config.Read(configPath)
	errFatal(err)

	switch len(os.Args) {
	case 2:
		fmt.Printf("Custom behaviour is set to :%v\n", configObj.CustomBehaviour)
	case 3:
		if os.Args[2] == "true" {
			configObj.CustomBehaviour = true
		} else if os.Args[2] == "false" {
			configObj.CustomBehaviour = false
		} else {
			fmt.Println("Invalid argument, expected 'true' or 'false'")
		}
		err = config.Write(configObj, configPath)
		errFatal(err)
	default:
		fmt.Println("Invalid number of arguments")
	}

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

    custom            Get custom behaviour
    custom bool       Set custom behaviour to true or false
                      (Custom behavior omits the path from the command execution,
                      running 'cmd' instead of 'cmd path')

    help              Print this help message

`)
}

func handleFuzzy() {
	p := fzf.StartFzf(configPath)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
