package main

import (
	"log"
	"net/http"
	"os"
	"time"

	customerv1 "github.com/umbranian0/customer-mdm/api/gen/api/proto/customer/v1"
	"github.com/umbranian0/customer-mdm/internal/adapters/transport/rest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	grpcTarget := envOr("GRPC_ADDR", "127.0.0.1:8080")
	restAddr := envOr("REST_ADDR", ":8090")

	conn, err := grpc.Dial(grpcTarget, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	handler := &rest.Gateway{
		Client:       customerv1.NewCustomerServiceClient(conn),
		CallTimeout:  10 * time.Second,
		MaxPageSize:  200,
		DefaultLimit: 50,
	}

	srv := &http.Server{
		Addr:         restAddr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("REST gateway listening on %s -> gRPC %s\n", restAddr, grpcTarget)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
