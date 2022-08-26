package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Urls []string `yaml:"urls"`
}

func (c *Config) Load() *Config {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(homedir)
	yamlFile, err := ioutil.ReadFile(homedir + "/.ssh/authorized_keys.yml")
	if err != nil {
		log.Printf("ioutil read %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("yaml unmarshal: %v", err)
	}
	return c
}
