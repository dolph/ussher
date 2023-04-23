package main

import (
	"log"
	"os"
)

func initLog() {
	file, err1 := os.OpenFile("/var/log/ussher/ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err1 != nil {
		file, err2 := os.OpenFile("ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		log.SetOutput(file)
		log.Print(err1)
		if err2 != nil {
			log.Print(err2)
		}
	}

	log.SetOutput(file)
}

func main() {
	initLog()

	// Refuse to run as root
	if runningAsRoot() {
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
