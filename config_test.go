package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createTempConfig(t *testing.T, content string) (string, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "testuser.yml")
	if err := os.WriteFile(tmpFile, []byte(content), 0o644); err != nil {
		t.Fatal("Failed to create temp config file:", err)
	}
	return tmpFile, func() {}
}

func TestConfigLoad(t *testing.T) {
	// Test case: valid YAML content
	validYamlContent := `
sources:
- url: https://example.com/keys
`
	tmpFile, cleanup := createTempConfig(t, validYamlContent)
	defer cleanup()

	config := &Config{}
	tmpFilePath, err := filepath.Abs(tmpFile)
	if err != nil {
		t.Errorf("Failed to get the absolute path of tmpFile: %v", err)
	}
	config.LoadConfigByPath(tmpFilePath)

	expectedSources := []Source{
		{URL: "https://example.com/keys"},
	}

	if len(config.Sources) != len(expectedSources) {
		t.Errorf("Expected %d sources, got %d", len(expectedSources), len(config.Sources))
	} else {
		for i, expected := range expectedSources {
			if config.Sources[i] != expected {
				t.Errorf("Expected source %d to be %v, got %v", i, expected, config.Sources[i])
			}
		}
	}

	// Test case: invalid YAML content using a dict of sources instead of a
	// list
	invalidYamlContent := `
---
sources:
  url: https://example.com/keys
`
	tmpFile, cleanup = createTempConfig(t, invalidYamlContent)
	defer cleanup()
	tmpFilePath, err = filepath.Abs(tmpFile)
	if err != nil {
		t.Errorf("Failed to get the absolute path of tmpFile: %v", err)
	}

	config = &Config{}
	config.LoadConfigByPath(tmpFilePath)

	if len(config.Sources) != 0 {
		t.Errorf("Expected 0 sources, got %d: %v", len(config.Sources), config.Sources[0])
	}
}

func TestResolveCacheTTL(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want time.Duration
	}{
		{"empty falls back to default", "", defaultCacheTTL},
		{"unparseable falls back to default", "five minutes", defaultCacheTTL},
		{"valid seconds", "30s", 30 * time.Second},
		{"valid minutes", "10m", 10 * time.Minute},
		{"valid hours", "1h", time.Hour},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Config{CacheTTL: tc.raw}
			if got := c.ResolveCacheTTL(); got != tc.want {
				t.Errorf("ResolveCacheTTL() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestConfigLoadCacheTTL(t *testing.T) {
	yamlContent := `
cache_ttl: 15m
sources:
- url: https://example.com/keys
`
	tmpFile, cleanup := createTempConfig(t, yamlContent)
	defer cleanup()

	config := &Config{}
	config.LoadConfigByPath(tmpFile)

	if config.CacheTTL != "15m" {
		t.Errorf("Expected CacheTTL %q, got %q", "15m", config.CacheTTL)
	}
	if got := config.ResolveCacheTTL(); got != 15*time.Minute {
		t.Errorf("ResolveCacheTTL() = %v, want %v", got, 15*time.Minute)
	}
}

func TestResolveHTTPTimeout(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want time.Duration
	}{
		{"empty falls back to default", "", defaultHTTPTimeout},
		{"unparseable falls back to default", "ten seconds", defaultHTTPTimeout},
		{"valid milliseconds", "500ms", 500 * time.Millisecond},
		{"valid seconds", "30s", 30 * time.Second},
		{"valid minutes", "2m", 2 * time.Minute},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := &Config{HTTPTimeout: tc.raw}
			if got := c.ResolveHTTPTimeout(); got != tc.want {
				t.Errorf("ResolveHTTPTimeout() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestConfigLoadHTTPTimeout(t *testing.T) {
	yamlContent := `
http_timeout: 5s
sources:
- url: https://example.com/keys
`
	tmpFile, cleanup := createTempConfig(t, yamlContent)
	defer cleanup()

	config := &Config{}
	config.LoadConfigByPath(tmpFile)

	if config.HTTPTimeout != "5s" {
		t.Errorf("Expected HTTPTimeout %q, got %q", "5s", config.HTTPTimeout)
	}
	if got := config.ResolveHTTPTimeout(); got != 5*time.Second {
		t.Errorf("ResolveHTTPTimeout() = %v, want %v", got, 5*time.Second)
	}
}

func TestConfigLoadRejectsWorldWritable(t *testing.T) {
	validYamlContent := `
sources:
- url: https://example.com/keys
`
	tmpFile, cleanup := createTempConfig(t, validYamlContent)
	defer cleanup()

	if err := os.Chmod(tmpFile, 0o666); err != nil {
		t.Fatal(err)
	}

	config := &Config{}
	config.LoadConfigByPath(tmpFile)

	if len(config.Sources) != 0 {
		t.Fatalf("expected no sources loaded from world-writable config, got %d", len(config.Sources))
	}
}
