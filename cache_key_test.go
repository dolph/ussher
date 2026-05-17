package main

import "testing"

func TestHTTPCacheKey(t *testing.T) {
	url := "https://git.example.com/users/alice/keys"

	if got := httpCacheKey(url, ""); got != url {
		t.Fatalf("httpCacheKey without credential = %q; want %q", got, url)
	}

	a := httpCacheKey(url, "token-a")
	b := httpCacheKey(url, "token-b")
	if a == b {
		t.Fatalf("expected different cache keys for different tokens, both %q", a)
	}
	if a == url || b == url {
		t.Fatal("expected credential to affect cache key")
	}
}
