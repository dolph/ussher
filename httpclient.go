// Package provides an HTTP client which uses a persistent cache. The caching
// behavior is not RFC 7234 compliant by design, and is not involved at the
// HTTP transport layer.
package main

import (
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	http  *http.Client
	cache *Cache
	ttl   time.Duration
}

func NewHTTPClient(ttl, httpTimeout time.Duration) *Client {
	return &Client{
		http: &http.Client{
			Timeout: httpTimeout,
		},
		cache: NewCache("/var/cache/ussher"),
		ttl:   ttl,
	}
}

// fresh reports whether a cached entry written at setAt is still within the
// configured TTL.
func (c *Client) fresh(setAt time.Time) bool {
	return time.Since(setAt) < c.ttl
}

func (c *Client) GetURL(url string) []string {
	if cached, setAt, ok := c.cache.Get(url); ok && c.fresh(setAt) {
		return bodyToKeys(cached)
	}

	log.Printf("GET %v", url)
	resp, err := c.http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch %v: %v", url, err)
		return make([]string, 0)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body from %v: %v", url, err)
			return make([]string, 0)
		}
		c.cache.Set(url, bodyBytes)
		return bodyToKeys(bodyBytes)
	}
	log.Printf("HTTP %d from %v", resp.StatusCode, url)
	return make([]string, 0)
}

func bodyToKeys(body []byte) []string {
	s := strings.TrimSuffix(string(body), "\n")
	if s == "" {
		log.Printf("Found %v key(s)", 0)
		return []string{}
	}
	parts := strings.Split(s, "\n")
	keys := make([]string, 0, len(parts))
	for _, line := range parts {
		if line != "" {
			keys = append(keys, line)
		}
	}
	log.Printf("Found %v key(s)", len(keys))
	return keys
}
