package structs

type Config struct {
	EditorCmd  string     `json:"editorCmd"`
	DirAliases []DirAlias `json:"aliases"`
}

type DirAlias struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
