package main

import (
	"log"
	"net"
	cf "courier/config"
	cp "courier/genproto/courier"
	"courier/service"
	"courier/storage/postgres"

	"google.golang.org/grpc"
)

func main() {
	config := cf.Load()

	db, err := postgres.NewPostgresStorage(config)

	if err != nil {
		panic(err)
	}

	listener, err := net.Listen("tcp", config.COURIER_SERVICE_PORT)

	if err != nil {
		log.Fatalf("Failed to listen tcp: %v", err)
	}

	s := grpc.NewServer()

	cp.RegisterNotificationServiceServer(s, service.NewNotificationService(db))
	cp.RegisterTaskServiceServer(s, service.NewTaskService(db))
	// op.RegisterCartItemServiceServer(s, service.NewCartItemService(db))
	// op.RegisterOrderServiceServer(s, service.NewOrderService(db))
	// op.RegisterOrderItemServiceServer(s, service.NewOrderItemService(db))

	log.Printf("server listening at %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
