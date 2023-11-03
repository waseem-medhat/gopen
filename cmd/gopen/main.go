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
	case "help", "h":
		handleHelp()

	case "init", "i":
		handleInit()

	case "editor", "e":
		handleEditor()

	case "alias", "a":
		handleAlias()

	default:
		handleGopen()
	}
}

func handleHelp() {
	const width = 16
	const strTmpl = "    %-*s  %s\n"

	fmt.Println("Gopen - a simple CLI to quick-start coding projects")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "gopen foo", "cd into path assigned to alias `foo` and run the editor cmd")
	fmt.Printf(strTmpl, width, "gopen cmd [args]", "Run command `cmd` (see Commands below)")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("Can be abbreviated by the first letter (`gopen i` == `gopen init`)")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "init", "Initialize a new config file (~/.config/gopen/gopen.json)")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "editor", "Get editor command")
	fmt.Printf(strTmpl, width, "editor cmd", "Set editor command to `cmd`")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "alias", "List all saved aliases")
	fmt.Printf(strTmpl, width, "alias foo", "Get path assigned to alias 'foo'")
	fmt.Printf(strTmpl, width, "alias foo bar", "Assign to alias `foo` the path `bar`")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "remove foo", "Remove alias `foo` from the config")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "help", "Print this help message")
	fmt.Println("")
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
