package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Urls []string `yaml:"urls"`
}

func (c *Config) Load() *Config {
	yamlFile, err := ioutil.ReadFile("/home/dolph/.ssh/authorized_keys.yml")
	if err != nil {
		log.Printf("ioutil read %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("yaml unmarshal: %v", err)
	}
	return c
}
