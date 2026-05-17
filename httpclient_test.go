package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// newTestClient builds a Client wired to a caller-supplied http.Client and a
// fresh temp cache directory. Returns the client plus a cleanup func.
func newTestClient(t *testing.T, httpClient *http.Client) (*Client, func()) {
	t.Helper()
	tmpDir := t.TempDir()
	c := &Client{
		http:  httpClient,
		cache: NewCache(tmpDir),
		ttl:   time.Minute,
	}
	return c, func() {}
}

func TestGetURL_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ssh-rsa AAAA alice\nssh-ed25519 BBBB bob\n")
	}))
	defer srv.Close()

	c, cleanup := newTestClient(t, srv.Client())
	defer cleanup()

	keys := c.GetURL(srv.URL)
	if len(keys) != 2 {
		t.Fatalf("want 2 keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "ssh-rsa AAAA alice" || keys[1] != "ssh-ed25519 BBBB bob" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestGetURL_5xxIsNotFatal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	c, cleanup := newTestClient(t, srv.Client())
	defer cleanup()

	keys := c.GetURL(srv.URL)
	if len(keys) != 0 {
		t.Errorf("want empty slice on 5xx, got %v", keys)
	}
}

func TestGetURL_BodyReadErrorIsNotFatal(t *testing.T) {
	// Hijack the connection and close it after writing only headers, so the
	// client gets an error while reading the body.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("server does not support hijacking")
		}
		conn, bw, err := hj.Hijack()
		if err != nil {
			t.Fatal("hijack:", err)
		}
		// 200 with a Content-Length larger than what we'll actually send,
		// then close the connection so io.ReadAll returns an unexpected EOF.
		fmt.Fprint(bw, "HTTP/1.1 200 OK\r\nContent-Length: 100\r\n\r\n")
		bw.Flush()
		conn.Close()
	}))
	defer srv.Close()

	c, cleanup := newTestClient(t, srv.Client())
	defer cleanup()

	keys := c.GetURL(srv.URL)
	if len(keys) != 0 {
		t.Errorf("want empty slice on body read error, got %v", keys)
	}
}

func TestGetURL_UnreachableIsNotFatal(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	// Close before the request — the URL points at a no-longer-listening port.
	srv.Close()

	c, cleanup := newTestClient(t, &http.Client{Timeout: 2 * time.Second})
	defer cleanup()

	keys := c.GetURL(url)
	if len(keys) != 0 {
		t.Errorf("want empty slice on unreachable host, got %v", keys)
	}
}

func TestGetURL_TimeoutIsNotFatal(t *testing.T) {
	// Server that accepts the connection but stalls before responding,
	// simulating a hung upstream. Without http.Client.Timeout the test would
	// hang for the OS-default; with a tight timeout the client must give up
	// promptly and GetURL must return empty (not panic, not log.Fatal).
	released := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-released
	}))
	defer func() {
		close(released)
		srv.Close()
	}()

	timeout := 100 * time.Millisecond
	c, cleanup := newTestClient(t, &http.Client{Timeout: timeout})
	defer cleanup()

	start := time.Now()
	keys := c.GetURL(srv.URL)
	elapsed := time.Since(start)

	if len(keys) != 0 {
		t.Errorf("want empty slice on timeout, got %v", keys)
	}
	if elapsed > 5*timeout {
		t.Errorf("expected timeout to fire near %s, took %s", timeout, elapsed)
	}
}

func TestBodyToKeys(t *testing.T) {
	tests := []struct {
		name string
		body string
		want []string
	}{
		{"empty body", "", []string{}},
		{"newline only", "\n", []string{}},
		{"single key", "ssh-rsa AAAA\n", []string{"ssh-rsa AAAA"}},
		{"skips blank lines", "a\n\nb\n", []string{"a", "b"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bodyToKeys([]byte(tc.body))
			if len(got) != len(tc.want) {
				t.Fatalf("len = %d, want %d: %v", len(got), len(tc.want), got)
			}
			for i := range tc.want {
				if got[i] != tc.want[i] {
					t.Errorf("got[%d] = %q, want %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

func TestNewHTTPClient_PropagatesTimeout(t *testing.T) {
	c := NewHTTPClient(time.Minute, 7*time.Second)
	if got := c.http.Timeout; got != 7*time.Second {
		t.Errorf("NewHTTPClient http.Timeout = %s, want 7s", got)
	}
}

