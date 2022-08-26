package main

func main() {
	var c Config
	c.Load()
	Run(&c)
}
