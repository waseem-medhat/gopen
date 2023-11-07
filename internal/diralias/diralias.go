// Package diralias contains functions for listing or modifying directory
// aliases in a Gopen config.
package diralias

import (
	"fmt"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// List pretty-prints each alias and its corresponding path
func List(config structs.Config) (fmtAliases []string) {
	var width int
	for _, dirAlias := range config.DirAliases {
		if len(dirAlias.Alias) > width {
			width = len(dirAlias.Alias)
		}
	}

	for _, dirAlias := range config.DirAliases {
		fmtAlias := fmt.Sprintf("%*s: %s", width, dirAlias.Alias, dirAlias.Path)
		fmtAliases = append(fmtAliases, fmtAlias)
	}

	return
}

// Add takes a config, a new alias, and its path, then it returns a new config
// struct with the newly added alias. If the alias already exists, the function
// will overwrite it. It also ensures that no alias matches Gopen commands like
// `alias` or `init`.
func Add(config structs.Config, alias string, path string) (newConfig structs.Config, err error) {
	newConfig = config

	// Check if alias is reserved
	// Prevents aliases from matching Gopen commands
	reserved := []string{"a", "alias", "e", "editor", "h", "help", "i", "init", "custom", "c", "m", "migrate", "r", "remove"}
	for _, r := range reserved {
		if r == alias {
			err = fmt.Errorf("Error: `%v` is reserved and can't be used as an alias", alias)
			return
		}
	}

	newDirAlias := structs.DirAlias{Alias: alias, Path: path}

	for i, dirAlias := range newConfig.DirAliases {
		if dirAlias.Alias == alias {
			newConfig.DirAliases[i] = newDirAlias
			return
		}
	}

	newConfig.DirAliases = append(newConfig.DirAliases, newDirAlias)
	return
}
