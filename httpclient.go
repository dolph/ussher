// Package provides an HTTP client which uses a persistent cache. The caching
// behavior is not RFC 7234 compliant by design, and is not involved at the
// HTTP transport layer.
package main

import (
	"encoding/json"
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

type GHEKey struct {
	ID  int    `json:"id"`
	Key string `json:"key"`
}

func (c *Client) GetGHE(ghe GithubEnterprise) []string {
	url := "https://" + ghe.Hostname + "/users/" + ghe.Username + "/keys"

	/* This will cache responses regardless of the token in context here, which could be a security risk. */
	if cached, setAt, ok := c.cache.Get(url); ok && c.fresh(setAt) {
		var keys []GHEKey
		err := json.Unmarshal(cached, &keys)
		if err != nil {
			log.Printf("Failed to decode JSON from cache for %v: %v", url, err)
			return make([]string, 0)
		}
		return GHEKeysToKeys(keys)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to build request for %v: %v", url, err)
		return make([]string, 0)
	}

	req.Header = http.Header{
		"Accept":        {"application/vnd.github+json"},
		"Authorization": {"token " + ghe.Token}}

	log.Printf("GET %v", url)
	resp, err := c.http.Do(req)
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
		var keys []GHEKey
		if err := json.Unmarshal(bodyBytes, &keys); err != nil {
			log.Printf("Failed to decode JSON from %v: %v", url, err)
			return make([]string, 0)
		}

		return GHEKeysToKeys(keys)
	}
	log.Printf("HTTP %d from %v", resp.StatusCode, url)
	return make([]string, 0)
}

func bodyToKeys(body []byte) []string {
	s := string(body)
	s = strings.TrimSuffix(s, "\n")
	keys := strings.Split(s, "\n")
	log.Printf("Found %v key(s)", len(keys))
	return keys
}

func GHEKeysToKeys(gheKeys []GHEKey) []string {
	var keys []string
	for _, gheKey := range gheKeys {
		keys = append(keys, gheKey.Key)
	}
	log.Printf("Found %v key(s)", len(keys))
	return keys
}
