package main

import (
	"strings"
	"testing"
)

func TestValidateArgs(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{"no args (just program name)", []string{"ussher"}, "", true},
		{"empty args slice", []string{}, "", true},
		{"single arg looks like username", []string{"ussher", "alice"}, "alice", false},
		{"single arg is --version", []string{"ussher", "--version"}, "--version", false},
		{"too many args", []string{"ussher", "alice", "extra"}, "", true},
		{"too many args with --version", []string{"ussher", "--version", "extra"}, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := validateArgs(tc.args)
			if (err != nil) != tc.wantErr {
				t.Fatalf("validateArgs(%v) err = %v, wantErr %v", tc.args, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("validateArgs(%v) = %q, want %q", tc.args, got, tc.want)
			}
		})
	}
}

// captureRunMain runs runMain with a deps that records which side-effect
// hooks fired and returns those flags alongside the error.
type runMainResult struct {
	err           error
	initLogCalled bool
	loadCalled    bool
	loadUsername  string
	runCalled     bool
}

func captureRunMain(args []string, execWritable, asRoot, validUser bool) runMainResult {
	var r runMainResult
	d := runDeps{
		executableWritable: func() bool { return execWritable },
		runningAsRoot:      func() bool { return asRoot },
		initLog:            func() { r.initLogCalled = true },
		validUser:          func(string) bool { return validUser },
		loadConfig: func(u string, _ *Config) {
			r.loadCalled = true
			r.loadUsername = u
		},
		runFn: func(*Config) { r.runCalled = true },
	}
	r.err = runMain(args, d)
	return r
}

func TestRunMain(t *testing.T) {
	cases := []struct {
		name             string
		args             []string
		execWritable     bool
		asRoot           bool
		validUser        bool
		wantErrSubstring string // empty = expect nil
		wantInitLog      bool
		wantRun          bool
	}{
		{
			name:             "missing args",
			args:             []string{"ussher"},
			validUser:        true,
			wantErrSubstring: "usage:",
		},
		{
			name: "version short-circuits before any gate",
			args: []string{"ussher", "--version"},
			// gates flipped to "would fail" — none should fire.
			execWritable: true,
			asRoot:       true,
			validUser:    false,
		},
		{
			name:             "executable writable",
			args:             []string{"ussher", "alice"},
			execWritable:     true,
			validUser:        true,
			wantErrSubstring: "permissions issue",
		},
		{
			name:             "running as root",
			args:             []string{"ussher", "alice"},
			asRoot:           true,
			validUser:        true,
			wantErrSubstring: "as root",
		},
		{
			name:             "invalid user",
			args:             []string{"ussher", "alice"},
			validUser:        false,
			wantErrSubstring: "User not found",
			wantInitLog:      true,
		},
		{
			name:        "happy path",
			args:        []string{"ussher", "alice"},
			validUser:   true,
			wantInitLog: true,
			wantRun:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := captureRunMain(tc.args, tc.execWritable, tc.asRoot, tc.validUser)

			if (r.err != nil) != (tc.wantErrSubstring != "") {
				t.Fatalf("err = %v, want substring %q", r.err, tc.wantErrSubstring)
			}
			if r.err != nil && !strings.Contains(r.err.Error(), tc.wantErrSubstring) {
				t.Errorf("err = %q, want substring %q", r.err.Error(), tc.wantErrSubstring)
			}
			if r.initLogCalled != tc.wantInitLog {
				t.Errorf("initLog called = %v, want %v", r.initLogCalled, tc.wantInitLog)
			}
			if r.runCalled != tc.wantRun {
				t.Errorf("runFn called = %v, want %v", r.runCalled, tc.wantRun)
			}
			if r.runCalled {
				if !r.loadCalled || r.loadUsername != "alice" {
					t.Errorf("loadConfig called with %q, want alice", r.loadUsername)
				}
			}
		})
	}
}

func TestDefaultDeps_AllSet(t *testing.T) {
	d := defaultDeps()
	if d.executableWritable == nil ||
		d.runningAsRoot == nil ||
		d.initLog == nil ||
		d.validUser == nil ||
		d.loadConfig == nil ||
		d.runFn == nil {
		t.Error("defaultDeps left a hook unset")
	}
}
