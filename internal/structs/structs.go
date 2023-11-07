// Package structs simply contains struct types to be used by other packages.
package structs

// Config is the struct for the v2 config file
type Config struct {
	EditorCmd       string     `json:"editorCmd"`
	CustomBehaviour bool       `json:"customBehaviour"`
	DirAliases      []DirAlias `json:"aliases"`
}

// ConfigV1 is the struct for the v1 config file
// It is used for migrating v1 config files to v2 config files and testing
type ConfigV1 struct {
	EditorCmd  string     `json:"editorCmd"`
	DirAliases []DirAlias `json:"aliases"`
}
type DirAlias struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
