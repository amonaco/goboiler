package main

import (
	"log"

	"github.com/amonaco/goboiler/api"
	"github.com/pelletier/go-toml"
)

const configFile = "./config.toml"

func main() {

	config, err := toml.LoadFile(configFile)
	if err != nil {
		log.Fatal("Error loading configuration file.")
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
