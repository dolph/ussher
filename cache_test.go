package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCache(t *testing.T) {
	// Create a temporary directory for the cache
	tempDir, err := ioutil.TempDir("", "cache-test")
	if err != nil {
		t.Fatal("Failed to create temporary directory for cache:", err)
	}
	defer os.RemoveAll(tempDir)

	cache := NewCache(tempDir)

	testKey := "test_key"
	testValue := []byte("test_value")

	// Test the Get method on an empty cache
	value, ok := cache.Get(testKey)
	if ok {
		t.Error("Expected cache miss, got cache hit")
	}

	// Test the Set method
	cache.Set(testKey, testValue)

	// Test the Get method after setting a value
	value, ok = cache.Get(testKey)
	if !ok {
		t.Error("Expected cache hit, got cache miss")
	}
	if !bytes.Equal(value, testValue) {
		t.Errorf("Expected value %q, got %q", testValue, value)
	}

	// Test the Delete method
	cache.Delete(testKey)

	// Test the Get method after deleting a value
	value, ok = cache.Get(testKey)
	if ok {
		t.Error("Expected cache miss, got cache hit")
	}
}
