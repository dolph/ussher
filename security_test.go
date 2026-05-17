package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsFileWorldWritable(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keys")
	if err := os.WriteFile(path, []byte("ssh-rsa AAAA"), 0o644); err != nil {
		t.Fatal(err)
	}

	writable, err := isFileWorldWritable(path)
	if err != nil {
		t.Fatalf("isFileWorldWritable: %v", err)
	}
	if writable {
		t.Error("expected 0644 file to not be world-writable")
	}

	if err := os.Chmod(path, 0o666); err != nil {
		t.Fatal(err)
	}
	writable, err = isFileWorldWritable(path)
	if err != nil {
		t.Fatalf("isFileWorldWritable after chmod: %v", err)
	}
	if !writable {
		t.Error("expected 0666 file to be world-writable")
	}
}
