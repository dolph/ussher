// Package cache provides a caching implementation that uses the diskv package
// to supplement an in-memory map with persistent storage.
package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/peterbourgon/diskv"
)

type Cache struct {
	d *diskv.Diskv
}

func NewCache(basePath string) *Cache {
	return &Cache{
		d: diskv.New(diskv.Options{
			BasePath:     basePath,
			CacheSizeMax: 100 * 1024 * 1024, // 100MB
		}),
	}
}

// Get returns the cached value for key along with the time it was written.
// The caller is responsible for deciding whether setAt is recent enough to
// trust. ok is false on a cache miss or on any read/parse error.
func (c *Cache) Get(key string) (value []byte, setAt time.Time, ok bool) {
	filename := keyToFilename(key)
	if filename == "" {
		log.Printf("Skipping unusable cache: %v", filename)
		return nil, time.Time{}, false
	}
	raw, err := c.d.Read(filename)
	if err != nil {
		log.Print("Cache MISS: ", key)
		return nil, time.Time{}, false
	}
	nl := bytes.IndexByte(raw, '\n')
	if nl < 0 {
		log.Print("Cache MISS (no header): ", key)
		return nil, time.Time{}, false
	}
	ts, err := strconv.ParseInt(string(raw[:nl]), 10, 64)
	if err != nil {
		log.Print("Cache MISS (bad header): ", key)
		return nil, time.Time{}, false
	}
	log.Print("Cache HIT: ", key)
	return raw[nl+1:], time.Unix(ts, 0), true
}

func (c *Cache) Set(key string, value []byte) {
	filename := keyToFilename(key)
	// Prepend a unix-timestamp header line so Get can report write time and
	// callers can enforce their own freshness policy.
	header := []byte(strconv.FormatInt(time.Now().Unix(), 10) + "\n")
	payload := append(header, value...)
	if err := c.d.WriteStream(filename, bytes.NewReader(payload), true); err != nil {
		log.Printf("Failed to write %v to cache: %v", key, err)
		return
	}
	log.Print("Cache SET: ", key)
}

func (c *Cache) Delete(key string) {
	filename := keyToFilename(key)
	if err := c.d.Erase(filename); err != nil {
		log.Printf("Failed to delete %v from cache: %v", key, err)
		return
	}
	log.Print("Cache DELETE: ", key)
}

func keyToFilename(key string) string {
	hash := sha1.New()
	if _, err := io.WriteString(hash, key); err != nil {
		log.Printf("Failed to generate cache filename from key: %v", err)
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}
