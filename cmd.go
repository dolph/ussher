package main

import (
	"fmt"
	"sync"
)

func Run(c *Config) {
	var wg sync.WaitGroup

	client := NewHTTPClient()
	keyChan := make(chan []string)
	for _, source := range c.Sources {
		wg.Add(1)
		go func(source Source) {
			defer wg.Done()
			if source.URL != "" {
				keyChan <- client.GetURL(source.URL)
			}
		}(source)
	}

	// Close keyChan whenever it's done
	go func() {
		wg.Wait()
		close(keyChan)
	}()

	// Output keys from keyChan
	for results := range keyChan {
		for _, k := range results {
			fmt.Println(k)
		}
	}
}
