package main

import (
	"github.com/dolph/ussher/pkg/cmd"
	"github.com/dolph/ussher/pkg/config"
)

func main() {
	var c config.Config
	c.Load()
	cmd.Run(&c)
}
