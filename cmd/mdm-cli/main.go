package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v5/pgxpool"
)

func main() {
    fmt.Println("mdm-cli available commands: migrate (auto-applies SQL in migrations/)")
    // This skeleton keeps CLI minimal; migrations run on service startup automatically.
    // Extend here if you want a dedicated migration runner.
    // Example:
    // pool, _ := pgxpool.New(context.Background(), "postgres://mdm:mdm@localhost:5432/mdm?sslmode=disable")
    // defer pool.Close()
    // _ = pool // use it to run SQL scripts
    log.Println("Done.")
    _ = context.Background()
    _ = pgxpool.Pool{}
}
