package fzf

import (
	"fmt"

	"github.com/wipdev-tech/gopen/internal/config"
)

func calcMaxWidths(aliases []config.DirAlias) (int, int, int) {
	maxAliasW := 0
	maxPathW := 0

	for _, a := range aliases {
		if len(a.Alias) > maxAliasW {
			maxAliasW = len(a.Alias)
		}
		if len(a.Path) > maxPathW {
			maxPathW = len(a.Path)
		}
	}

	return maxAliasW, maxPathW, maxAliasW + maxPathW + 6
}

func alignQuestion(question string, maxW int) string {
	fmtStr := fmt.Sprintf("%%-%ds", maxW)
	return fmt.Sprintf(fmtStr, question)
}

func alignResult(alias, path string, maxAliasW, maxPathW int) string {
	fmtStr := fmt.Sprintf("  %%-%ds  %%-%ds ", maxAliasW, maxPathW+1)
	return fmt.Sprintf(fmtStr, alias, path)
}
