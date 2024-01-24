package config_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/wipdev-tech/gopen/internal/config"
)

func TestInitConfigCreatesNewFile(t *testing.T) {
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	configPath := dir + "/config.json"
	err = config.Init(dir, configPath)
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

	err = config.Init(dir, configPath)
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
	err = config.Init(dir, configPath)
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

	err = config.Init(dir, configPath)
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
	err = config.Init(dir, configPath)
	if err != nil {
		t.Fatal(err)
	}

	newConfig, err := config.Read(configPath)
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
	_, err := config.Read("/tmp/nonexistent_file")
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
	cfg, err := config.Read(tmpfile2.Name())
	if err != nil {
		t.Fatal(err)
	}
	expected := config.C{
		EditorCmd: "vim",
		DirAliases: []config.DirAlias{
			{Alias: "docs", Path: "/usr/share/doc"},
		},
	}
	if !reflect.DeepEqual(cfg, expected) {
		t.Fatalf("Expected %v but got %v", expected, cfg)
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
	_, err = config.Read(tmpfile3.Name())
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

	testConfig := config.C{
		EditorCmd: "vim",
		DirAliases: []config.DirAlias{
			{Alias: "docs", Path: "/usr/share/doc"},
		},
	}
	expectedOutput := `{
  "editorCmd": "vim",
  "customBehaviour": false,
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

	err = config.Write(testConfig, tmpfile.Name())
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
