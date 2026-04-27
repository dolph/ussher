package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	defaultCacheTTL     = 5 * time.Minute
	defaultHTTPTimeout  = 10 * time.Second
)

type GithubEnterprise struct {
	Hostname string `yaml:"api_hostname"`
	Username string `yaml:"user"`
	Token    string `yaml:"token"`
}

type Source struct {
	URL string `yaml:"url"`
}

type Config struct {
	// CacheTTL is parsed by time.ParseDuration ("5m", "30s", "1h"). An empty
	// or unparseable value falls back to defaultCacheTTL.
	CacheTTL string `yaml:"cache_ttl"`
	// HTTPTimeout caps how long a single upstream fetch (connect + headers +
	// body read) can take before it's abandoned. Parsed by time.ParseDuration.
	// An empty or unparseable value falls back to defaultHTTPTimeout.
	HTTPTimeout string   `yaml:"http_timeout"`
	Sources     []Source `yaml:"sources"`
}

// ResolveCacheTTL returns the configured cache TTL, falling back to a sane
// default when unset or malformed.
func (c *Config) ResolveCacheTTL() time.Duration {
	if c.CacheTTL == "" {
		return defaultCacheTTL
	}
	d, err := time.ParseDuration(c.CacheTTL)
	if err != nil {
		log.Printf("Invalid cache_ttl %q, falling back to %s: %v", c.CacheTTL, defaultCacheTTL, err)
		return defaultCacheTTL
	}
	return d
}

// ResolveHTTPTimeout returns the configured per-request HTTP timeout,
// falling back to a sane default when unset or malformed.
func (c *Config) ResolveHTTPTimeout() time.Duration {
	if c.HTTPTimeout == "" {
		return defaultHTTPTimeout
	}
	d, err := time.ParseDuration(c.HTTPTimeout)
	if err != nil {
		log.Printf("Invalid http_timeout %q, falling back to %s: %v", c.HTTPTimeout, defaultHTTPTimeout, err)
		return defaultHTTPTimeout
	}
	return d
}

func (c *Config) LoadConfigByUser(username string) {
	// `username` is validated at this point to be a valid Linux username, so
	// it's safe to load this configuration file without the risk of loading
	// arbitrary paths.
	c.LoadConfigByPath("/etc/ussher/" + username + ".yml")
}

func (c *Config) LoadConfigByPath(path string) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to %v ", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Failed to parse as YAML: %v", err)
		return
	}
	log.Printf("Loaded configuration from %v", path)
}
