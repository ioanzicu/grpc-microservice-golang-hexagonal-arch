package main

// Here are handled dependency injections

import (
	"fmt"
	"log"

	"github.com/ioanzicu/microservices/order/config"
	"github.com/ioanzicu/microservices/order/internal/adapters/db"
	"github.com/ioanzicu/microservices/order/internal/adapters/grpc"
	"github.com/ioanzicu/microservices/order/internal/adapters/payment"
	"github.com/ioanzicu/microservices/order/internal/application/core/api"
)

func main() {
	dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	paymentAdapter, err := payment.NewAdapter(config.GetPaymentServiceURL())
	if err != nil {
		log.Fatalf("Failed to initialize payment stub. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter, paymentAdapter)
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())

	fmt.Println("Start the server...")
	grpcAdapter.Run()
}
