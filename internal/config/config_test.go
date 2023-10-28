package config

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/wipdev-tech/gopen/internal/structs"
)

func TestReadConfig(t *testing.T) {
	// Case 1: reading a file that does not exist
	_, err := ReadConfig("/tmp/nonexistent_file")
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
	config, err := ReadConfig(tmpfile2.Name())
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
	_, err = ReadConfig(tmpfile3.Name())
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

	testConfig := structs.Config{
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

	err = WriteConfig(testConfig, tmpfile.Name())
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
