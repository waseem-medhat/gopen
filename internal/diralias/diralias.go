// Package diralias contains functions for listing or modifying directory
// aliases in a Gopen config.
package diralias

import (
	"fmt"
	"path/filepath"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// listDirAliases pretty-prints each alias and its corresponding path
func List(config structs.Config) []string {
	var width int

	for _, dirAlias := range config.DirAliases {
		if len(dirAlias.Alias) > width {
			width = len(dirAlias.Alias)
		}
	}

    var fmtAliases []string
	for _, dirAlias := range config.DirAliases {
		fmtAlias := fmt.Sprintf("%*s: %s", width, dirAlias.Alias, dirAlias.Path)
		fmtAliases = append(fmtAliases, fmtAlias)
	}

	return fmtAliases
}

// Add takes a config, a new alias, and its path, then it returns a new config
// struct with the newly added alias. If the alias already exists, the function
// will overwrite it. It also ensures that no alias matches Gopen commands like
// `alias` or `init`.
func Add(config structs.Config, alias string, path string) (structs.Config, error) {
    newConfig := config

	reserved := []string{"a", "alias", "e", "editor", "h", "help", "i", "init"}
	for _, r := range reserved {
		if r == alias {
            err := fmt.Errorf("Error: `%v` is reserved and can't be used as an alias", alias)
			return newConfig, err
		}
	}

	// If the path is ".", then we want to use the current directory
	// instead of the literal "."
	if path == "." {
		path = "./"
	}
	newPath, err := filepath.Abs(path)
	if err != nil {
		return newConfig, err
	}

	newDirAlias := structs.DirAlias{Alias: alias, Path: newPath}

	for i, dirAlias := range newConfig.DirAliases {
		if dirAlias.Alias == alias {
			newConfig.DirAliases[i] = newDirAlias
			return newConfig, err
		}
	}

	newConfig.DirAliases = append(newConfig.DirAliases, newDirAlias)
	return newConfig, err
}
