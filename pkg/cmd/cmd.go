package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/dolph/ssh-doorman/pkg/config"
)

func getURL(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return ""
		}
		return string(bodyBytes)
	}
	return ""
}

func Run(c *config.Config) {
	var wg sync.WaitGroup

	keyChan := make(chan string)

	for _, url := range c.Urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			keyChan <- getURL(url)
		}(url)
	}

	// Close keyChan whenever it's done
	go func() {
		wg.Wait()
		close(keyChan)
	}()

	// Collect keys from keyChan
	var keys []string
	for k := range keyChan {
		keys = append(keys, k)
		fmt.Print(k)
	}
}
