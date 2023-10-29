package diralias

import (
	"testing"

	"github.com/wipdev-tech/gopen/internal/structs"
)

func TestListDirAliases(t *testing.T) {
	config := structs.Config{DirAliases: []structs.DirAlias{}}
	result := ListDirAliases(config)
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

	actual := ListDirAliases(config)

	for i, actLine := range actual {
		if actLine != expected[i] {
			t.Fatalf("Got\n%v but expected\n%v\n", actLine, expected[i])
		}
	}
}
