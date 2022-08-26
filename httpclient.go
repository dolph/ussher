// Package provides an HTTP client which uses a persistent cache. The caching
// behavior is not RFC 7234 compliant by design, and is not involved at the
// HTTP transport layer.
package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client struct {
	http  *http.Client
	cache *Cache
}

func NewHTTPClient() *Client {
	return &Client{
		http:  &http.Client{},
		cache: NewCache("/var/cache/ussher"),
	}
}

func (c *Client) GetURL(url string) []string {
	if cached, ok := c.cache.Get(url); ok {
		return bodyToKeys(cached)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return make([]string, 0)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
			return make([]string, 0)
		}
		c.cache.Set(url, bodyBytes)
		return bodyToKeys(bodyBytes)
	}
	return make([]string, 0)
}

func bodyToKeys(body []byte) []string {
	s := string(body)
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
