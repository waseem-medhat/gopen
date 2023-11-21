// Package gopen includes the main Gopen execution function.
package gopen

import (
	"errors"
	"os"
	"os/exec"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// Gopen uses the Config struct to find the path corresponding to targetAlias
// and executes the editor command with the target path as the working
// directory
func Gopen(targetAlias string, config structs.Config) error {
	var targetPath string
	for _, dirAlias := range config.DirAliases {
		if targetAlias == dirAlias.Alias {
			targetPath = dirAlias.Path
			break
		}
	}

	if targetPath == "" {
		return errors.New("Invalid command or non-existent alias\nRun `gopen help` for info")
	}

	editorCmd := config.EditorCmd
    err := os.Chdir(targetPath)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	// CustomBehaviour lets the user open the target path in a new buffer
	if config.CustomBehaviour {
		cmd = exec.Command(editorCmd)
	} else {
		cmd = exec.Command(editorCmd, targetPath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}
