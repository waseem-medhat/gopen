// Package diralias contains functions for listing or modifying directory
// aliases in a Gopen config.
package diralias

import (
	"fmt"
	"path/filepath"

	"github.com/wipdev-tech/gopen/internal/config"
)

// List pretty-prints each alias and its corresponding path
func List(config config.C) []string {
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
func Add(cfg config.C, alias string, path string) (config.C, error) {
	newCfg := cfg

	reserved := []string{"a", "alias", "e", "editor", "h", "help", "i", "init"}
	for _, r := range reserved {
		if r == alias {
			err := fmt.Errorf("Error: `%v` is reserved and can't be used as an alias", alias)
			return newCfg, err
		}
	}

	// If the path is ".", then we want to use the current directory
	// instead of the literal "."
	if path == "." {
		path = "./"
	}
	newPath, err := filepath.Abs(path)
	if err != nil {
		return newCfg, err
	}

	newDirAlias := config.DirAlias{Alias: alias, Path: newPath}

	for i, dirAlias := range cfg.DirAliases {
		if dirAlias.Alias == alias {
			newCfg.DirAliases[i] = newDirAlias
			return newCfg, err
		}
	}

	newCfg.DirAliases = append(newCfg.DirAliases, newDirAlias)
	return newCfg, err
}
