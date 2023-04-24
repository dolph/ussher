package main

import (
	"os"
	"os/user"
	"regexp"
)

// Return true if ussher is running as the root user, which would violate
// the principle of least-privilege.
func isRunningAsRoot() bool {
	return os.Getuid() == 0
}

// Ensure that the input string is a valid Linux account name on this host.
// This prevents security issues such as:
// - Reading arbitrary files on the host
// - Log injection
func isValidLinuxAccountName(name string) bool {
	// Check if the input string matches the allowed character pattern
	var validNamePattern = regexp.MustCompile("^[a-z_][a-z0-9_-]*$")
	if !validNamePattern.MatchString(name) {
		return false
	}

	// Check if the input string is within the allowed length
	if len(name) > 32 {
		return false
	}

	// Check if the input string is already an existing user account on the host
	_, err := user.Lookup(name)
	if err == nil {
		return false
	}

	return true
}

// Ensures we're not reading a file that can be easily modified by an attacker.
func isFileWorldWritable(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		// If we can't stat the file, it either doesn't exist or just assume
		// it's world writeable (failsafe)
		return true, err
	}

	permissions := fileInfo.Mode().Perm()

	// Check if the world writable bit is set (i.e., 0002)
	if permissions&os.ModePerm&0o0002 != 0 {
		return true, nil
	}

	return false, nil
}
