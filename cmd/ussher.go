package main

import (
	"github.com/dolph/ssh-doorman/pkg/cmd"
	"github.com/dolph/ssh-doorman/pkg/config"
)

func main() {
	var c config.Config
	c.Load()
	cmd.Run(&c)
}
