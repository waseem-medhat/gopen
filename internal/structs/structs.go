// Package structs simply contains struct types to be used by other packages.
package structs

type Config struct {
	EditorCmd       string     `json:"editorCmd"`
	CustomBehaviour bool       `json:"customBehaviour"`
	DirAliases      []DirAlias `json:"aliases"`
}

type DirAlias struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
