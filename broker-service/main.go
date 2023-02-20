package main

import (
	"log"

	"github.com/gitsuki/finance/broker/api"
)

func main() {
	server, err := api.NewServer()
	if err != nil {
		log.Fatal("[fatal] cannot create server", err)
	}
}
