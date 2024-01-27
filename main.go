// This is the main entry point for the Gopen tool. It includes the parsing
// logic for command-line arguments and handler functions associated with each
// command.
package main

import (
	"os"

	"github.com/wipdev-tech/gopen/internal/handlers"
)

var configDir = os.Getenv("HOME") + "/.config/gopen"
var configPath = configDir + "/gopen.json"

func main() {
	if len(os.Args) < 2 {
		// handleFzf()
		handlers.Fzf(configDir, configPath)
		return
	}

	switch os.Args[1] {
	case "help", "h":
		handlers.Help(configDir, configPath)

	case "init", "i":
		handlers.Init(configDir, configPath)

	case "editor", "e":
		handlers.Editor(configDir, configPath)

	case "alias", "a":
		handlers.Alias(configDir, configPath)

	case "remove", "r":
		handlers.Remove(configDir, configPath)

	case "custom", "c":
		handlers.Custom(configDir, configPath)

	default:
		handlers.Gopen(configDir, configPath)
	}
}
