// Package cache provides a caching implementation that uses the diskv package
// to supplement an in-memory map with persistent storage.
package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"

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

func (c *Cache) Get(key string) (value []byte, ok bool) {
	filename := keyToFilename(key)
	value, err := c.d.Read(filename)
	if err != nil {
		log.Print("Cache MISS: ", key)
		return []byte{}, false
	}
	log.Print("Cache HIT: ", key)
	return value, true
}

func (c *Cache) Set(key string, value []byte) {
	filename := keyToFilename(key)
	err := c.d.WriteStream(filename, bytes.NewReader(value), true)
	if err != nil {
		log.Printf("Failed to write %v to cache:", key, err)
		return
	}
	log.Print("Cache SET: ", key)
}

func (c *Cache) Delete(key string) {
	filename := keyToFilename(key)
	err := c.d.Erase(filename)
	if err != nil {
		log.Printf("Failed to delete %v from cache:", key, err)
		return
	}
	log.Print("Cache DELETE: ", key)
}

func keyToFilename(key string) string {
	hash := sha1.New()
	io.WriteString(hash, key)
	return hex.EncodeToString(hash.Sum(nil))
}
