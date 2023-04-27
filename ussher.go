package main

import (
	"log"
	"os"
)

func main() {
	// Security sanity checks
	if isRunningWorldWritable() {
		log.Fatal("Refusing to run world writable binary")
	}
	if isRunningAsRoot() {
		log.Fatal("Refusing to run as root")
	}

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
	c.Load(username)
	Run(&c)
}
