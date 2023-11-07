package config

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/wipdev-tech/gopen/internal/structs"
)

func TestInitConfigCreatesNewFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	configPath := dir + "/config.json"
	err = Init(dir, configPath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(configPath)
	if os.IsNotExist(err) {
		t.Errorf("expected file to exist at %s, but it doesn't", configPath)
	}
}

func TestInitConfigReturnsErrorIfFileExists(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	configPath := dir + "/config.json"
	_, err = os.Create(configPath)
	if err != nil {
		t.Fatal(err)
	}

	err = Init(dir, configPath)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
	if !os.IsExist(err) {
		t.Errorf("expected error to be os.ErrExist, but got %v", err)
	}
}

func TestInitConfigReturnsErrorIfDirectoryCreationFails(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	configPath := dir + "/nonexistent/config.json"
	err = Init(dir, configPath)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func TestInitConfigReturnsErrorIfFileCreationFails(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	configPath := dir + "/config.json"
	file, err := os.Create(configPath)
	if err != nil {
		t.Fatal(err)
	}

	err = file.Chmod(0400)
	if err != nil {
		t.Fatal(err)
	}

	err = Init(dir, configPath)
	if err == nil {
		t.Error("expected an error, but got nil")
	}
}

func TestInitConfigWritesEmptyConfig(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	configPath := dir + "/config.json"
	err = Init(dir, configPath)
	if err != nil {
		t.Fatal(err)
	}

	newConfig, err := Read(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if newConfig.EditorCmd != "" {
		t.Fatal("Config's EditorCmd is not empty")
	}

	if len(newConfig.DirAliases) != 0 {
		t.Fatal("Config's DirAliases is not empty")
	}
}

func TestReadConfig(t *testing.T) {
	// Case 1: reading a file that does not exist
	_, err := Read("/tmp/nonexistent_file")
	if !os.IsNotExist(err) {
		t.Fatalf("Expected a \"Not exist\" error but got \"%v\"", err)
	}

	// Case 2: reading a valid config file
	tmpfile2, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile2.Name())
	_, err = tmpfile2.Write([]byte(`{"editorCmd": "vim", "aliases": [{"alias": "docs", "path": "/usr/share/doc"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpfile2.Close(); err != nil {
		t.Fatal(err)
	}
	config, err := Read(tmpfile2.Name())
	if err != nil {
		t.Fatal(err)
	}
	expected := structs.Config{
		EditorCmd: "vim",
		DirAliases: []structs.DirAlias{
			{Alias: "docs", Path: "/usr/share/doc"},
		},
	}
	if !reflect.DeepEqual(config, expected) {
		t.Fatalf("Expected %v but got %v", expected, config)
	}

	// Case 3: reading an invalid JSON
	tmpfile3, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile3.Name())
	_, err = tmpfile3.Write([]byte("l{\"editor\": \"vim\", \"aliases\": [{\"name\": \"docs\", \"path\": \"/usr/share/doc\"}]}"))
	if err != nil {
		t.Fatal(err)
	}
	if err := tmpfile3.Close(); err != nil {
		t.Fatal(err)
	}
	_, err = Read(tmpfile3.Name())
	if err == nil {
		t.Fatal("Expected an error but got nil")
	}
	if _, ok := err.(*json.SyntaxError); !ok {
		t.Fatalf("Expected a *json.SyntaxError but got %T", err)
	}

}

func TestWriteConfig(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	testConfig := structs.ConfigV1{
		EditorCmd: "vim",
		DirAliases: []structs.DirAlias{
			{Alias: "docs", Path: "/usr/share/doc"},
		},
	}
	expectedOutput := `{
  "editorCmd": "vim",
  "aliases": [
    {
      "alias": "docs",
      "path": "/usr/share/doc"
    }
  ]
}`
	err = os.WriteFile(tmpfile.Name(), []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = WriteV1(testConfig, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	fileContents, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if string(fileContents) != expectedOutput {
		t.Fatalf("Expected %q but got %q", expectedOutput, string(fileContents))
	}
}

func TestReadFromOldConfig(t *testing.T) {
	var oldConfig structs.ConfigV1
	oldConfig.EditorCmd = "vim"
	oldConfig.DirAliases = []structs.DirAlias{
		{Alias: "docs", Path: "/usr/share/doc"},
	}
	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = WriteV1(oldConfig, tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	newConfig, err := Read(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if newConfig.CustomBehaviour != false {
		t.Fatal("Config's CustomBehaviour is not false")
	}

}

func TestMigrate(t *testing.T) {
	var oldConfig structs.ConfigV1
	oldConfig.EditorCmd = "vim"
	oldConfig.DirAliases = []structs.DirAlias{
		{Alias: "docs", Path: "/usr/share/doc"},
	}

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	err = WriteV1(oldConfig, tmpfile.Name())

	err = Migrate(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	newConfig, err := Read(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if newConfig.CustomBehaviour != false {
		t.Fatal("Config's CustomBehaviour is not false")
	}

	if newConfig.EditorCmd != "vim" {
		t.Fatal("Config's EditorCmd is not vim")
	}

	if len(newConfig.DirAliases) != 1 {
		t.Fatal("Config's DirAliases is not 1")
	}

	if newConfig.DirAliases[0].Alias != "docs" {
		t.Fatal("Config's DirAliases[0].Alias is not docs")
	}

	if newConfig.DirAliases[0].Path != "/usr/share/doc" {
		t.Fatal("Config's DirAliases[0].Path is not /usr/share/doc")
	}
}
