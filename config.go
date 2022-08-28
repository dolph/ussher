package main

import (
	"io/ioutil"
	"log"

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

func (c *Config) Load(username string) *Config {
	path := "/home/" + username + "/.ssh/authorized_keys.yml"
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("ioutil read %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("yaml unmarshal: %v", err)
	}
	log.Printf("Loaded configuration from %v", path)
	return c
}
