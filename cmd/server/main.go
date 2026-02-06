package main

import (
	"log"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := server.NewHTTPServer(conf)
	if err := srv.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
