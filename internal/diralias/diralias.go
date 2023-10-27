package diralias

import (
	"fmt"

	"github.com/wipdev-tech/gopen/internal/structs"
)

// listDirAliases pretty-prints each alias and its corresponding path
func ListDirAliases(config structs.Config) {
	var width int
	for _, dirAlias := range config.DirAliases {
		if len(dirAlias.Alias) > width {
			width = len(dirAlias.Alias)
		}
	}

	for _, dirAlias := range config.DirAliases {
		fmt.Printf("%*s: %s\n", width, dirAlias.Alias, dirAlias.Path)
	}
}
