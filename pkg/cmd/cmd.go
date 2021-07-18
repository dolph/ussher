package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/dolph/ssh-doorman/pkg/config"
)

func getURL(url string, wg *sync.WaitGroup) string {
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
	return ""
}

func Run(c *config.Config) {
	var wg sync.WaitGroup

	for _, url := range c.Urls {
		wg.Add(1)
		go getURL(url, &wg)
	}
	wg.Wait()
}
