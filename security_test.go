package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidUser(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		expectedValue bool
	}{
		{
			name:          "Valid user (root)",
			username:      "root",
			expectedValue: true,
		},
		{
			name:          "Valid user (nobody)",
			username:      "nobody",
			expectedValue: true,
		},
		{
			name:          "Invalid user - starts with number",
			username:      "1kofabhhfsbf6krb",
			expectedValue: false,
		},
		{
			name:          "Invalid user - contains uppercase letters",
			username:      "XD5hObIMZF2zKS7W",
			expectedValue: false,
		},
		{
			name:          "Invalid user - too long",
			username:      "idfkjcacexia1dyji5iwcfweoliamzpn1",
			expectedValue: false,
		},
		{
			name:          "Invalid user - contains invalid characters",
			username:      "moctcg!@",
			expectedValue: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isValidUser(test.username)
			if result != test.expectedValue {
				t.Errorf("Expected %v, got %v", test.expectedValue, result)
			}
		})
	}
}

func TestUidIsRoot(t *testing.T) {
	cases := []struct {
		uid  int
		want bool
	}{
		{0, true},
		{1, false},
		{500, false},
		{1000, false},
		{-1, false},
	}
	for _, tc := range cases {
		if got := uidIsRoot(tc.uid); got != tc.want {
			t.Errorf("uidIsRoot(%d) = %v, want %v", tc.uid, got, tc.want)
		}
	}
}

// writeFileMode creates path with the exact mode requested, defeating umask.
func writeFileMode(t *testing.T, path string, mode os.FileMode) {
	t.Helper()
	if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
}

func TestIsPathUnsafelyWritable(t *testing.T) {
	dir := t.TempDir()

	cases := []struct {
		name string
		mode os.FileMode
		want bool
	}{
		{"safe 0755", 0755, false},
		{"safe 0750", 0750, false},
		{"safe 0700", 0700, false},
		{"group writable 0775", 0775, true},
		{"group writable 0770", 0770, true},
		{"world writable 0757", 0757, true},
		{"world writable 0707", 0707, true},
		{"both 0777", 0777, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(dir, tc.name)
			writeFileMode(t, path, tc.mode)
			if got := isPathUnsafelyWritable(path); got != tc.want {
				t.Errorf("isPathUnsafelyWritable(mode=%o) = %v, want %v", tc.mode, got, tc.want)
			}
		})
	}

	t.Run("missing path fails safe to true", func(t *testing.T) {
		if !isPathUnsafelyWritable(filepath.Join(dir, "does-not-exist")) {
			t.Error("missing path should be reported as unsafely writable (failsafe)")
		}
	})
}

func TestIsFileWorldWritable(t *testing.T) {
	dir := t.TempDir()

	t.Run("non-world-writable", func(t *testing.T) {
		path := filepath.Join(dir, "safe")
		writeFileMode(t, path, 0640)
		got, err := isFileWorldWritable(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got {
			t.Errorf("isFileWorldWritable(0640) = true, want false")
		}
	})

	t.Run("world-writable", func(t *testing.T) {
		path := filepath.Join(dir, "permissive")
		writeFileMode(t, path, 0666)
		got, err := isFileWorldWritable(path)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !got {
			t.Errorf("isFileWorldWritable(0666) = false, want true")
		}
	})

	t.Run("missing file returns true and an error", func(t *testing.T) {
		got, err := isFileWorldWritable(filepath.Join(dir, "missing"))
		if err == nil {
			t.Error("expected error for missing file")
		}
		if !got {
			t.Error("missing file should be reported as world-writable (failsafe)")
		}
	})
}
