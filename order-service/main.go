package main

import (
	"log"
	cf "order/config"
	op "order/genproto/order"
	"order/storage/postgres"
	"order/service"
	"net"

	"google.golang.org/grpc"
)

func main() {
	config := cf.Load()

	db, err := postgres.NewPostgresStorage(config)

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", config.ORDER_SERVICE_PORT)

	if err != nil {
		log.Fatalf("Failed to listen tcp: %v", err)
	}

	s := grpc.NewServer()

	op.RegisterProductServiceServer(s, service.NewProductService(db))
	op.RegisterCartServiceServer(s, service.NewCartService(db))
	// op.RegisterMemoriesServiceServer(s, service.NewMemoriesService(db))
	// op.RegisterSharedMemoriesServiceServer(s, service.NewSharedMemoriesService(db))

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}