package main

import (
	"log"

	"github.com/noilpa/gobalance/internal/config"
	"github.com/noilpa/gobalance/internal/server"
)

func main() {
	path := "config/test/config.yml"
	cfg, err := config.LoadConfig(path)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	s, err := server.New(cfg.ListenPort, cfg.Strategy, cfg.Backends)
	if err != nil {
		log.Fatalf("failed to create new server: %v", err)
	}

	if err := s.Start(); err != nil {
		log.Fatalf("failed to create new server: %v", err)
	}
}
