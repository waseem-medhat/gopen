// Package diralias contains functions for listing or modifying directory
// aliases in a Gopen config.
package diralias

import (
	"fmt"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// listDirAliases pretty-prints each alias and its corresponding path
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

// Add takes a config, a new alias, and its path, then it returns a new config object
func Add(config structs.Config, alias string, path string) (newConfig structs.Config) {
	newConfig = config
	newConfig.DirAliases = append(
		newConfig.DirAliases,
		structs.DirAlias{Alias: alias, Path: path},
	)

	return
}
