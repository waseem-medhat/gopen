package handlers

import (
	"fmt"
	"os"

	"github.com/wipdev-tech/gopen/internal/config"
	"github.com/wipdev-tech/gopen/internal/fzf"
)

func Init(configDir, configPath string) {
	err := config.Init(configDir, configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}
}

func Editor(configDir, configPath string) {
	cfg, err := config.Read(configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		return
	}

	if len(os.Args) < 3 {
		fmt.Println(cfg.EditorCmd)
		return
	}

	cfg.EditorCmd = os.Args[2]
	err = config.Write(cfg, configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}

}

func Alias(configDir, configPath string) {
	cfg, err := config.Read(configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		return
	}

	switch len(os.Args) {
	case 2:
		for _, fmtAlias := range cfg.ListAliases() {
			fmt.Println(fmtAlias)
		}

	case 3:
		for _, dirAlias := range cfg.DirAliases {
			if dirAlias.Alias == os.Args[2] {
				fmt.Println(dirAlias.Path)
				return
			}
		}
		fmt.Println("Alias doesn't exist")

	case 4:
		cfg, err := cfg.AddAlias(os.Args[2], os.Args[3])
		if err != nil {
			fmt.Println(err)
			return
		}

		err = config.Write(cfg, configPath)
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
		}

	default:
		fmt.Println("Too many arguments - exiting...")
	}
}

func Gopen(configDir, configPath string) {
	cfg, err := config.Read(configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		return
	}

	err = cfg.Gopen(os.Args[1])
	if err != nil {
		fmt.Println(err)
	}
}

func Remove(configDir, configPath string) {
	cfg, err := config.Read(configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		return
	}

	if len(os.Args) != 3 {
		fmt.Println("Error: must provide one alias to 'remove' command")
		return
	}

	var newConfig config.C
	newConfig.EditorCmd = cfg.EditorCmd
	for _, dirAlias := range cfg.DirAliases {
		if dirAlias.Alias != os.Args[2] {
			newConfig.DirAliases = append(newConfig.DirAliases, dirAlias)
		}
	}

	err = config.Write(newConfig, configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
	}
}

func Custom(configDir, configPath string) {
	cfg, err := config.Read(configPath)
	if err != nil {
		fmt.Println(fmt.Errorf("error: %v", err))
		return
	}

	switch len(os.Args) {
	case 2:
		fmt.Printf("Custom behaviour is set to :%v\n", cfg.CustomBehaviour)
	case 3:
		if os.Args[2] == "true" {
			cfg.CustomBehaviour = true
		} else if os.Args[2] == "false" {
			cfg.CustomBehaviour = false
		} else {
			fmt.Println("Invalid argument, expected 'true' or 'false'")
		}
		err = config.Write(cfg, configPath)
		if err != nil {
			fmt.Println(fmt.Errorf("error: %v", err))
		}
	default:
		fmt.Println("Invalid number of arguments")
	}

}

func Help(configDir, configPath string) {
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

func Fzf(configDir, configPath string) {
	p := fzf.StartFzf(configPath)
	m, err := p.Run()
	if err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if fzfModel, ok := m.(fzf.Model); ok {
		alias := fzfModel.Selected
		if alias != "" {
			err = fzfModel.Config.Gopen(alias)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
