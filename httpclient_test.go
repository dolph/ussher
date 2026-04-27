package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// newTestClient builds a Client wired to a caller-supplied http.Client and a
// fresh temp cache directory. Returns the client plus a cleanup func.
func newTestClient(t *testing.T, httpClient *http.Client) (*Client, func()) {
	t.Helper()
	tmpDir, err := ioutil.TempDir("", "httpclient-test")
	if err != nil {
		t.Fatal("temp dir:", err)
	}
	c := &Client{
		http:  httpClient,
		cache: NewCache(tmpDir),
		ttl:   time.Minute,
	}
	return c, func() { os.RemoveAll(tmpDir) }
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

func TestGetGHE_UnreachableIsNotFatal(t *testing.T) {
	c, cleanup := newTestClient(t, &http.Client{Timeout: 2 * time.Second})
	defer cleanup()

	// example.invalid is reserved and guaranteed to not resolve, exercising
	// GetGHE's transport-error path. Pre-fix this would log.Fatal and kill
	// the test process.
	keys := c.GetGHE(GithubEnterprise{
		Hostname: "example.invalid",
		Username: "alice",
		Token:    "irrelevant",
	})
	if len(keys) != 0 {
		t.Errorf("want empty slice on unreachable GHE host, got %v", keys)
	}
}
