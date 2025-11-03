package main

import (
    "context"
    "log"
    "net"
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

    lis, err := net.Listen("tcp", ":8080")
    if err != nil { log.Fatal(err) }
    log.Println("gRPC listening on :8080")
    if err := c.GRPCServer.Serve(lis); err != nil {
        log.Fatal(err)
    }
    time.Sleep(time.Second)
}
