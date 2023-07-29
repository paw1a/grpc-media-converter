package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/paw1a/grpc-media-converter/api_gateway_service/config"
	"github.com/paw1a/grpc-media-converter/api_gateway_service/internal/client"
	delivery "github.com/paw1a/grpc-media-converter/api_gateway_service/internal/delivery/http"
	storage "github.com/paw1a/grpc-media-converter/api_gateway_service/pb/storage"
	_ "github.com/paw1a/grpc-media-converter/api_gateway_service/pkg/logging"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	storageServiceConn, err := client.NewStorageServiceConn(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create storage service connection: %v", err)
	}
	storageService := storage.NewStorageServiceClient(storageServiceConn)

	handler := delivery.NewHandler(cfg, storageService)
	server := &http.Server{
		Handler:      handler.Init(),
		Addr:         fmt.Sprintf(":%s", cfg.Http.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Infof("server started on port %s", cfg.Http.Port)
	log.Fatal(server.ListenAndServe())
}
