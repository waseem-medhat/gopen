// Package diralias contains functions for listing or modifying directory
// aliases in a Gopen config.
package diralias

import (
	"fmt"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// listDirAliases pretty-prints each alias and its corresponding path
func ListDirAliases(config structs.Config) []string {
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
