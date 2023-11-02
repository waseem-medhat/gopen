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
	const width = 17
	const strTmpl = "    %-*s  %s\n"

	fmt.Println("Gopen - a simple CLI to quick-start coding projects")
	fmt.Println("")
	fmt.Println(`
The premise of this command-line utility is to save an editor of choice and a
list of aliases for your local development projects instead of "polluting" your
system-level configs (e.g., .bashrc). Then, Gopen command will cd into that
folder and open your editor of choice.`,
	)
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("")
	fmt.Printf(strTmpl, width, "init, i", "Initialize a new config file")
	fmt.Printf(strTmpl, width, "editor, e", "Get editor command")
	fmt.Printf(strTmpl, width, "    editor cmd", "Set editor command to `cmd`")
	fmt.Printf(strTmpl, width, "alias, a", "List all saved aliases")
	fmt.Printf(strTmpl, width, "    alias foo", "Get path for alias 'foo'")
	fmt.Printf(strTmpl, width, "    alias foo bar", "Set alias `foo` to path `bar`")
	fmt.Printf(strTmpl, width, "remove, r alias", "Remove `alias` from the list")
	fmt.Printf(strTmpl, width, "help, h", "Print this help message")
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
