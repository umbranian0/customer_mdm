module github.com/stricker/customer-mdm

go 1.22

require (
    github.com/jackc/pgx/v5 v5.6.0
    github.com/segmentio/kafka-go v0.4.46
    google.golang.org/grpc v1.66.2
    google.golang.org/protobuf v1.33.0
)

replace github.com/stricker/customer-mdm => ./
