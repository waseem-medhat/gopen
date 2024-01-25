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

func TestList(t *testing.T) {
	cfg := config.C{DirAliases: []config.DirAlias{}}
	result := cfg.ListAliases()
	if len(result) != 0 {
		t.Errorf("Expected an empty slice, but got %v", result)
	}

	cfg = config.C{
		DirAliases: []config.DirAlias{
			{Alias: "x", Path: "/path/to/x"},
			{Alias: "yz", Path: "/path/to/yz"},
			{Alias: "abc", Path: "/path/to/abc"},
		},
	}

	expected := []string{
		"  x: /path/to/x",
		" yz: /path/to/yz",
		"abc: /path/to/abc",
	}

	actual := cfg.ListAliases()

	for i, actLine := range actual {
		if actLine != expected[i] {
			t.Fatalf("Got\n%v but expected\n%v\n", actLine, expected[i])
		}
	}
}

func TestAdd(t *testing.T) {
	// Test adding a new alias
	cfg := config.C{
		EditorCmd: "",
		DirAliases: []config.DirAlias{
			{Alias: "alias1", Path: "/path/to/dir1"},
			{Alias: "alias2", Path: "/path/to/dir2"},
		},
	}
	newConfig, err := cfg.AddAlias("alias3", "/path/to/dir3")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedConfig := config.C{
		EditorCmd: "",
		DirAliases: []config.DirAlias{
			{Alias: "alias1", Path: "/path/to/dir1"},
			{Alias: "alias2", Path: "/path/to/dir2"},
			{Alias: "alias3", Path: "/path/to/dir3"},
		},
	}
	if !reflect.DeepEqual(newConfig, expectedConfig) {
		t.Errorf("Expected %v, but got %v", expectedConfig, newConfig)
	}

	// Test overwriting an existing alias
	newConfig, err = cfg.AddAlias("alias2", "/path/to/newdir")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedConfig = config.C{
		EditorCmd: "",
		DirAliases: []config.DirAlias{
			{Alias: "alias1", Path: "/path/to/dir1"},
			{Alias: "alias2", Path: "/path/to/newdir"},
		},
	}
	if !reflect.DeepEqual(newConfig, expectedConfig) {
		t.Errorf("Expected %v, but got %v", expectedConfig, newConfig)
	}

	// Test adding a reserved alias
	_, err = cfg.AddAlias("alias", "/path/to/newdir")
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	expectedError := "Error: `alias` is reserved and can't be used as an alias"
	if err.Error() != expectedError {
		t.Errorf("Expected %q, but got %q", expectedError, err.Error())
	}
}
