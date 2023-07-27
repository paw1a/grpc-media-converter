package main

import (
	"flag"
	"github.com/paw1a/grpc-media-converter/storage_service/config"
	"github.com/paw1a/grpc-media-converter/storage_service/internal/server"
	"log"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := server.NewServer(cfg)
	log.Fatal(s.Run())
}
