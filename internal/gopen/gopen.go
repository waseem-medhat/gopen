// Package gopen includes the main Gopen execution function.
package gopen

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// Gopen uses the Config struct to find the path corresponding to targetAlias
// and executes the editor command with the target path as the working
// directory
func Gopen(targetAlias string, config structs.Config) {
	var targetPath string
	for _, dirAlias := range config.DirAliases {
		if targetAlias == dirAlias.Alias {
			targetPath = dirAlias.Path
			break
		}
	}

	fInfo, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		fmt.Println("Path doesn't exist")
		return
	} else if err != nil {
		log.Fatal(err)
	}

	if !fInfo.IsDir() {
		println("Not a directory")
		return
	}

	editorCmd := config.EditorCmd
	os.Chdir(targetPath)
	cmd := exec.Command(editorCmd)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
