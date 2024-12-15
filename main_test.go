package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRename(t *testing.T) {
	got := rename("test.md")
	want := filepath.Join(os.TempDir(), "test.html")
	if got != want {
		t.Fatalf("got: %v, want: %v\n", got, want)
	}
}
