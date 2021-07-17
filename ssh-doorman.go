package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"gopkg.in/yaml.v2"
)

func getURL(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bodyBytes))
	}

}

type config struct {
	Urls []string `yaml:"urls"`
}

func (c *config) loadConfig() *config {
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

func main() {
	var c config
	c.loadConfig()

	var wg sync.WaitGroup

	for _, url := range c.Urls {
		wg.Add(1)
		go getURL(url, &wg)
	}
	wg.Wait()
}
