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
func Gopen(targetAlias string, config structs.Config) (err error) {
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
    var cmd *exec.Cmd
    editorCmd := config.EditorCmd
    err = os.Chdir(targetPath)
    if err != nil {
        return
    }
    // Custom editor allows for the use of non-terminal based editors
    // that need the path to the project as an argument
    if config.CustomBehaviour {
        cmd = exec.Command(editorCmd, targetPath)
    } else {
        cmd = exec.Command(editorCmd)
    }

    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    err = cmd.Run()
    if err != nil {
        return
    }

    return
}
