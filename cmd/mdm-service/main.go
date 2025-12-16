package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/umbranian0/customer-mdm/internal/app"
)

func main() {
	ctx := context.Background()
	c := app.Initialize(ctx)

	go func() {
		if err := c.OutboxDispatcher.Run(ctx); err != nil {
			log.Println("outbox dispatcher stopped:", err)
		}
	}()

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gRPC listening on", addr)
	if err := c.GRPCServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
}
