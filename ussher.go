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

func initLog() {
	file, err := os.OpenFile("/var/log/ussher/ussher.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func main() {
	initLog()
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
