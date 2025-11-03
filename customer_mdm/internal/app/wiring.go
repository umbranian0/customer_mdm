package app

import (
    "context"
    "log"
    "os"
    "time"
    "path/filepath"
    "io/ioutil"

    "github.com/jackc/pgx/v5/pgxpool"
    "google.golang.org/grpc"

    customerv1 "github.com/umbranian0/customer-mdm/api/gen/customer/v1"
    pg "github.com/umbranian0/customer-mdm/internal/adapters/db/postgres"
    outbox "github.com/umbranian0/customer-mdm/internal/adapters/cdc/outbox"
    "github.com/umbranian0/customer-mdm/internal/adapters/transport/grpc"
    "github.com/umbranian0/customer-mdm/internal/usecase"
)

type Container struct {
    GRPCServer *grpc.Server
    CustomerServer *grpcadp.CustomerServer
    OutboxDispatcher *outbox.Dispatcher
}

type Config struct {
    DB_DSN string
    KafkaBrokers string
    OutboxTopic string
}

func loadConfig() Config {
    cfg := Config{
        DB_DSN: os.Getenv("DB_DSN"),
        KafkaBrokers: os.Getenv("KAFKA_BROKERS"),
        OutboxTopic: os.Getenv("OUTBOX_TOPIC"),
    }
    if cfg.DB_DSN == "" {
        cfg.DB_DSN = "postgres://mdm:mdm@localhost:5432/mdm?sslmode=disable"
    }
    if cfg.KafkaBrokers == "" { cfg.KafkaBrokers = "localhost:9094" }
    if cfg.OutboxTopic == "" { cfg.OutboxTopic = "mdm.customer.events.v1" }
    return cfg
}

func autoMigrate(ctx context.Context, pool *pgxpool.Pool) {
    dir := "migrations"
    files, err := ioutil.ReadDir(dir)
    if err != nil {
        log.Println("migrate read dir:", err)
        return
    }
    for _, f := range files {
        if f.IsDir() { continue }
        path := filepath.Join(dir, f.Name())
        b, err := ioutil.ReadFile(path)
        if err != nil { log.Println("read migration:", path, err); continue }
        if _, err := pool.Exec(ctx, string(b)); err != nil {
            log.Println("apply migration:", path, err)
        } else {
            log.Println("applied migration:", path)
        }
    }
}

func Initialize(ctx context.Context) *Container {
    cfg := loadConfig()

    pool, err := pgxpool.New(ctx, cfg.DB_DSN)
    if err != nil { log.Fatal(err) }

    autoMigrate(ctx, pool)

    repo := &pg.CustomerRepository{Pool: pool}
    txm := &pg.TxManager{Pool: pool}
    outboxWriter := &pg.OutboxWriter{Pool: pool}

    createUC := &usecase.CreateCustomer{Repo: repo, Tx: txm, Outbox: outboxWriter, Topic: cfg.OutboxTopic}
    getUC    := &usecase.GetCustomer{Repo: repo}
    updateUC := &usecase.UpdateCustomer{Repo: repo, Tx: txm, Outbox: outboxWriter, Topic: cfg.OutboxTopic}
    deleteUC := &usecase.DeleteCustomer{Repo: repo, Tx: txm, Outbox: outboxWriter, Topic: cfg.OutboxTopic}
    listUC   := &usecase.ListCustomers{Repo: repo}

    srv := &grpc.Server{}
    custSrv := &grpcadp.CustomerServer{
        CreateUC: createUC, GetUC: getUC, UpdateUC: updateUC, DeleteUC: deleteUC, ListUC: listUC,
    }
    customerv1.RegisterCustomerServiceServer(srv, custSrv)

    dispatcher := &outbox.Dispatcher{
        Pool: pool,
        Publisher: &grpcadp.DummyPublisher{}, // swap to real Kafka in cmd main
        BatchSize: 100,
        PollEvery: 2 * time.Second,
        Topic: cfg.OutboxTopic,
    }

    return &Container{GRPCServer: srv, CustomerServer: custSrv, OutboxDispatcher: dispatcher}
}
