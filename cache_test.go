package main

import (
	"bytes"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewCache(t.TempDir())

	testKey := "test_key"
	testValue := []byte("test_value")

	if _, _, ok := cache.Get(testKey); ok {
		t.Error("Expected cache miss, got cache hit")
	}

	before := time.Now().Add(-time.Second)
	cache.Set(testKey, testValue)
	after := time.Now().Add(time.Second)

	value, setAt, ok := cache.Get(testKey)
	if !ok {
		t.Fatal("Expected cache hit, got cache miss")
	}
	if !bytes.Equal(value, testValue) {
		t.Errorf("Expected value %q, got %q", testValue, value)
	}
	if setAt.Before(before) || setAt.After(after) {
		t.Errorf("setAt %v not in [%v, %v]", setAt, before, after)
	}

	cache.Delete(testKey)

	if _, _, ok := cache.Get(testKey); ok {
		t.Error("Expected cache miss after delete, got cache hit")
	}
}
