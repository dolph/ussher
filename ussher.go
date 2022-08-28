package main

import (
	"log"
	"os"
)

func securityCheck() {
	if os.Getuid() == 0 {
		log.Fatal("ussher refuses to run as root")
	}
}

func main() {
	securityCheck()
	if len(os.Args) != 2 {
		log.Fatal("usage: ussher <username>")
	}
	username := os.Args[1]
	log.Print("Sourcing authorized_keys for ", username)
	var c Config
	c.Load(username)
	Run(&c)
}
