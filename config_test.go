package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func createTempConfig(t *testing.T, content string) (string, func()) {
	tmpDir, err := ioutil.TempDir("", "config-test")
	if err != nil {
		t.Fatal("Failed to create temp dir for config test:", err)
	}

	tmpFile := filepath.Join(tmpDir, "testuser.yml")
	err = ioutil.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatal("Failed to create temp config file:", err)
	}

	return tmpFile, func() { os.RemoveAll(tmpDir) }
}

func TestConfigLoad(t *testing.T) {
	// Test case: valid YAML content
	validYamlContent := `
sources:
  - url: "https://example.com/keys"
  - github_enterprise:
      api_hostname: "github.example.com"
      user: "testuser"
      token: "testtoken"
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
		{
			GHE: GithubEnterprise{
				Hostname: "github.example.com",
				Username: "testuser",
				Token:    "testtoken",
			},
		},
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

	// Test case: invalid YAML content
	invalidYamlContent := `
sources:
  - url: "https://example.com/keys"
  - github_enterprise:
      api_hostname: "github.example.com"
      user: "testuser"
      token: "testtoken
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
		t.Errorf("Expected 0 sources, got %d", len(config.Sources))
	}
}
