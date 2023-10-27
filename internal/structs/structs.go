// Package structs simply contains struct types to be used by other packages.
package structs

type Config struct {
	EditorCmd  string     `json:"editorCmd"`
	DirAliases []DirAlias `json:"aliases"`
}

type DirAlias struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
