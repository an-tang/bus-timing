package main

import (
	"bus-timing/cmd"

	config "bus-timing/configuration"
)

func main() {
	cmd.RunServer()
}

func init() {
	loadConfig()
}

func loadConfig() {
	err := config.LoadConfig("./configuration")
	if err != nil {
		panic(err)
	}
}
