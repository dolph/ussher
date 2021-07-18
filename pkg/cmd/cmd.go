package cmd

import (
	"fmt"
	"sync"

	"github.com/dolph/ssh-doorman/pkg/config"
	"github.com/dolph/ssh-doorman/pkg/httpclient"
)

func Run(c *config.Config) {
	var wg sync.WaitGroup

	client := httpclient.New()
	keyChan := make(chan []string)
	for _, url := range c.Urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			keyChan <- client.GetURL(url)
		}(url)
	}

	// Close keyChan whenever it's done
	go func() {
		wg.Wait()
		close(keyChan)
	}()

	// Collect keys from keyChan
	var keys []string
	for results := range keyChan {
		for _, k := range results {
			keys = append(keys, k)
			fmt.Println(k)
		}
	}
}
