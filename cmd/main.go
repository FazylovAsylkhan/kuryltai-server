package main

import (
	"log"

	"github.com/FazylovAsylkhan/kuryltai-server/cmd/service"
	"github.com/FazylovAsylkhan/kuryltai-server/internal/config"
)


func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Get()
	httpServer, err := service.Init(cfg)
	if err != nil {
		return err
	}
	httpServer.Start()	
	return nil
}