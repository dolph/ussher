package main

import (
	"log"
	"os"
)

func initLog() {
	// Attempt to write a log file to a standard location
	file, err1 := os.OpenFile("/var/log/ussher/ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		// If we can't write to a standard location, then try to write to the current working directory
		file, err2 := os.OpenFile("ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			// If we can't produce real log files, then use stdout and abort
			log.Fatal("Refusing to run without being able to log to /var/log/ussher/ (", err1, ") or current working directory (", err2, ")")
		}

		// At least we can log to the current working directory
		log.SetOutput(file)
		log.Print("Failed to write to /var/log/ussher: ", err1)
	}

	log.SetOutput(file)
}

func main() {
	// Support `ussher --version` to print versioning info about the executable
	// before proceeding to security-hardening checks.
	if os.Args[1] == "--version" {
		PrintVersion()
		return
	}

	// Security sanity checks
	if isExecutableWritable() {
		log.Fatal("Refusing to run due to permissions issue on the ussher executable")
	}
	if isRunningAsRoot() {
		log.Fatal("Refusing to run as root")
	}

	// Initialize logging AFTER security checks to ensure we're writing logs as
	// a non-root user
	initLog()

	// Check if the input username is provided
	if len(os.Args) != 2 {
		log.Fatal("usage: ussher <username>")
	}

	// Check if the input username is valid
	username := os.Args[1]
	if !isValidUser(username) {
		log.Fatal("User not found")
	}

	// At this point, we know that the input username is valid and safe to use.
	log.Print("Sourcing authorized_keys for ", username)

	var c Config
	c.LoadConfigByUser(username)
	Run(&c)
}
