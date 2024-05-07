// Package config includes functions and types for initializing, reading, and
// writing Gopen config files.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
)

// C is the struct representation of Gopen config.
type C struct {
	EditorCmd       string     `json:"editorCmd"`
	CustomBehaviour bool       `json:"customBehaviour"`
	DirAliases      []DirAlias `json:"aliases"`
}

// DirAlias is the struct type for the directory aliases where each struct
// contains the alias and the path it corresponds to.
type DirAlias struct {
	Alias   string `json:"alias"`
	Path    string `json:"path"`
	GitRepo string `json:"git_repo"`
}

// Init checks if the config file exists in configPath. If not, creates an
// empty config file. configDir will also be created if it doesn't exist.
func Init(configDir string, configPath string) error {
	_, err := os.Stat(configPath)
	if err == nil {
		return os.ErrExist
	}

	err = os.MkdirAll(configDir, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = os.Create(configPath)
	if err != nil {
		return err
	}

	emptyConfig := C{}
	err = Write(emptyConfig, configPath)
	if err != nil {
		return err
	}

	return err
}

// Write writes config to configPath (will OVERWRITE if file already exists)
func Write(config C, configPath string) error {
	jsonFile, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, jsonFile, 0644)
	if err != nil {
		return err
	}

	return err
}

// Read reads the configPath file and returns a Config struct
func Read(configPath string) (C, error) {
	var config C

	f, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(f, &config)
	if err != nil {
		return config, err
	}

	return config, err
}

// ListAliases pretty-prints each alias and its corresponding path
func (cfg C) ListAliases() []string {
	var width int

	for _, dirAlias := range cfg.DirAliases {
		if len(dirAlias.Alias) > width {
			width = len(dirAlias.Alias)
		}
	}

	var fmtAliases []string
	for _, dirAlias := range cfg.DirAliases {
		fmtAlias := fmt.Sprintf("%*s: %s", width, dirAlias.Alias, dirAlias.Path)
		fmtAliases = append(fmtAliases, fmtAlias)
	}

	return fmtAliases
}

// AddAlias takes a config, a new alias, and its path, then it returns a new
// config struct with the newly added alias. If the alias already exists, the
// function will overwrite it. It also ensures that no alias matches Gopen
// commands like `alias` or `init`.
func (cfg C) AddAlias(alias string, path string) (C, error) {
	newCfg := cfg

	reserved := []string{"a", "alias", "e", "editor", "h", "help", "i", "init", "g", "git"}
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

	newDirAlias := DirAlias{Alias: alias, Path: newPath}

	for i, dirAlias := range cfg.DirAliases {
		if dirAlias.Alias == alias {
			newCfg.DirAliases[i] = newDirAlias
			return newCfg, err
		}
	}

	newCfg.DirAliases = append(newCfg.DirAliases, newDirAlias)
	return newCfg, err
}

func (cfg C) SetGitRepo(alias string, repo string) (C, error) {
	for i, dirAlias := range cfg.DirAliases {
		if dirAlias.Alias == alias {
			newDirAlias := DirAlias{
				Alias:   dirAlias.Alias,
				Path:    dirAlias.Path,
				GitRepo: repo,
			}

			cfg.DirAliases[i] = newDirAlias
			return cfg, nil
		}
	}

	return cfg, fmt.Errorf("alias doesn't exist")
}

// Gopen uses the Config struct to find the path corresponding to targetAlias
// and executes the editor command with the target path as the working
// directory
func (cfg C) Gopen(targetAlias string) error {
	var targetPath string
	var targetRepo string
	for _, dirAlias := range cfg.DirAliases {
		if targetAlias == dirAlias.Alias {
			targetPath = dirAlias.Path
			targetRepo = dirAlias.GitRepo
			break
		}
	}

	if targetPath == "" {
		return errors.New("Invalid command or non-existent alias\nRun `gopen help` for info")
	}

	editorCmd := strings.Split(cfg.EditorCmd, " ")

	_, err := os.Stat(targetPath)
	if os.IsNotExist(err) && targetRepo != "" {
		fmt.Printf("dir %v not found\ntrying to clone %v\n", targetPath, targetRepo)
		_, err = git.PlainClone(targetPath, false, &git.CloneOptions{
			URL:      targetRepo,
			Progress: os.Stdout,
		})
	}
	if err != nil {
		return err
	}

	err = os.Chdir(targetPath)
	if err != nil {
		return err
	}

	var cmd *exec.Cmd
	// CustomBehaviour lets the user open the target path in a new buffer
	if cfg.CustomBehaviour {
		cmd = exec.Command(editorCmd[0], editorCmd[1:]...)
	} else {
		cmd = exec.Command(editorCmd[0], targetPath)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		return err
	}

	return err
}
