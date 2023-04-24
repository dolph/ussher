package main

import (
	"log"
	"os"
)

func main() {
	// Refuse to run as root
	if isRunningAsRoot() {
		log.Fatal("Refusing to run as root")
	}

	// Check if the input username is provided
	if len(os.Args) != 2 {
		log.Fatal("usage: ussher <username>")
	}

	// Check if the input username is valid
	username := os.Args[1]
	if !isValidLinuxAccountName(username) {
		log.Fatal("Invalid username")
	}

	// At this point, we know that the input username is valid and safe to use.
	log.Print("Sourcing authorized_keys for ", username)

	var c Config
	c.Load(username)
	Run(&c)
}
