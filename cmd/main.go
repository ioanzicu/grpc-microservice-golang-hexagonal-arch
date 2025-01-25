package main

import (
	"fmt"
	"log"

	"github.com/ioanzicu/microservices/order/config"
	"github.com/ioanzicu/microservices/order/internal/adapters/db"
	"github.com/ioanzicu/microservices/order/internal/adapters/grpc"
	"github.com/ioanzicu/microservices/order/internal/application/core/api"
)

func main() {
	fmt.Println("Start the server...")

	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}
