package main

import (
	"fmt"
	"os"
	"os/user"
	"regexp"
)

// Return true if ussher's binary is unnecessarily writable.
func isExecutableWritable() bool {
	executablePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Failed to get a path to ussher executable: %v\n", err)
		return true
	}
	return isPathUnsafelyWritable(executablePath)
}

// isPathUnsafelyWritable returns true if path's mode has the group-write
// or other-write bit set, or if path can't be stat'd (failsafe). Extracted
// from isExecutableWritable so the bit-mask logic is testable against an
// arbitrary chmod'd file rather than /proc/self/exe.
func isPathUnsafelyWritable(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Failed to stat %s: %v\n", path, err)
		return true
	}

	mode := fileInfo.Mode()

	// Check for group writable
	if mode&0020 != 0 {
		fmt.Printf("%s is group writable\n", path)
		return true
	}

	// Check for world writable
	if mode&0002 != 0 {
		fmt.Printf("%s is world writable\n", path)
		return true
	}

	return false
}

// Return true if ussher is running as the root user, which would violate
// the principle of least-privilege.
func isRunningAsRoot() bool {
	return uidIsRoot(os.Getuid())
}

// uidIsRoot returns true iff the given uid is 0. Extracted so the
// "is root?" predicate can be tested without actually running as root.
func uidIsRoot(uid int) bool {
	return uid == 0
}

// Ensure that the input string is a valid Linux account name on this host.
// This prevents security issues such as:
// - Reading arbitrary files on the host
// - Log injection
func isValidUser(name string) bool {
	// Check if the input string is within the allowed length
	if len(name) > 32 {
		return false
	}

	// Check if the input string matches the allowed character pattern
	var validNamePattern = regexp.MustCompile("^[a-z_][a-z0-9_-]*$")
	if !validNamePattern.MatchString(name) {
		return false
	}

	// Check if the input string is already an existing user account on the host
	_, err := user.Lookup(name)
	if err != nil {
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
	if permissions&os.ModePerm&0002 != 0 {
		return true, nil
	}

	return false, nil
}
