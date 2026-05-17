package main

import (
	"crypto/sha256"
	"encoding/hex"
)

// httpCacheKey returns a cache key for url, mixing in credential when present
// so different auth tokens cannot share cached responses (#4).
func httpCacheKey(url, credential string) string {
	if credential == "" {
		return url
	}
	sum := sha256.Sum256([]byte(credential))
	return url + "#auth=" + hex.EncodeToString(sum[:8])
}
