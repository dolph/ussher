package main

import "testing"

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
