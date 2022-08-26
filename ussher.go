package main

import (
	"github.com/dolph/ussher/cmd"
	"github.com/dolph/ussher/config"
)

func main() {
	var c config.Config
	c.Load()
	cmd.Run(&c)
}
