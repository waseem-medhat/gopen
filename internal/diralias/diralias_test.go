package diralias

import (
	"reflect"
	"testing"

	"github.com/wipdev-tech/gopen/internal/structs"
)

func TestList(t *testing.T) {
	config := structs.Config{DirAliases: []structs.DirAlias{}}
	result := List(config)
	if len(result) != 0 {
		t.Errorf("Expected an empty slice, but got %v", result)
	}

	config = structs.Config{
		DirAliases: []structs.DirAlias{
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

	actual := List(config)

	for i, actLine := range actual {
		if actLine != expected[i] {
			t.Fatalf("Got\n%v but expected\n%v\n", actLine, expected[i])
		}
	}
}

func TestAdd(t *testing.T) {
	// Test adding a new alias
	config := structs.Config{
		EditorCmd: "",
		DirAliases: []structs.DirAlias{
			{Alias: "alias1", Path: "/path/to/dir1"},
			{Alias: "alias2", Path: "/path/to/dir2"},
		},
	}
	newConfig, err := Add(config, "alias3", "/path/to/dir3")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedConfig := structs.Config{
		EditorCmd: "",
		DirAliases: []structs.DirAlias{
			{Alias: "alias1", Path: "/path/to/dir1"},
			{Alias: "alias2", Path: "/path/to/dir2"},
			{Alias: "alias3", Path: "/path/to/dir3"},
		},
	}
	if !reflect.DeepEqual(newConfig, expectedConfig) {
		t.Errorf("Expected %v, but got %v", expectedConfig, newConfig)
	}

	// Test overwriting an existing alias
	newConfig, err = Add(config, "alias2", "/path/to/newdir")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	expectedConfig = structs.Config{
		EditorCmd: "",
		DirAliases: []structs.DirAlias{
			{Alias: "alias1", Path: "/path/to/dir1"},
			{Alias: "alias2", Path: "/path/to/newdir"},
		},
	}
	if !reflect.DeepEqual(newConfig, expectedConfig) {
		t.Errorf("Expected %v, but got %v", expectedConfig, newConfig)
	}

	// Test adding a reserved alias
	_, err = Add(config, "alias", "/path/to/newdir")
	if err == nil {
		t.Errorf("Expected an error, but got nil")
	}
	expectedError := "Error: `alias` is reserved and can't be used as an alias"
	if err.Error() != expectedError {
		t.Errorf("Expected %q, but got %q", expectedError, err.Error())
	}
}
