package main

import (
	"log"

	"github.com/gitsuki/finance/broker/api"
	"github.com/gitsuki/finance/broker/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("[fatal] cannot load config", err)
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("[fatal] cannot create server", err)
	}
}
