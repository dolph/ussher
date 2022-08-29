package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type GithubEnterprise struct {
	Hostname string `yaml:"api_hostname"`
	Username string `yaml:"user"`
	Token    string `yaml:"token"`
}

type Source struct {
	URL string           `yaml:"url"`
	GHE GithubEnterprise `yaml:"github_enterprise"`
}

type Config struct {
	Sources []Source `yaml:"sources"`
}

func (c *Config) Load(username string) {
	path := "/home/" + username + "/.ssh/authorized_keys.yml"
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to %v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Failed to parse as YAML: %v", err)
		return
	}
	log.Printf("Loaded configuration from %v", path)
}
