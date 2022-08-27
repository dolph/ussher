package main

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Urls []string `yaml:"urls"`
}

func (c *Config) Load() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	yamlFile, err := ioutil.ReadFile(homeDir + "/.ssh/authorized_keys.yml")
	if err != nil {
		log.Printf("ioutil read %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("yaml unmarshal: %v", err)
	}
	return c
}
