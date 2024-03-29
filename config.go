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
	URL string `yaml:"url"`
}

type Config struct {
	Sources []Source `yaml:"sources"`
}

func (c *Config) LoadConfigByUser(username string) {
	// `username` is validated at this point to be a valid Linux username, so
	// it's safe to load this configuration file without the risk of loading
	// arbitrary paths.
	c.LoadConfigByPath("/etc/ussher/" + username + ".yml")
}

func (c *Config) LoadConfigByPath(path string) {
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
