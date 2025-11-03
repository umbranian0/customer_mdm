# Build
FROM golang:1.22 as builder
WORKDIR /src
COPY . .
RUN go mod download
# You can generate protobufs here if buf is available in the image; omitted for brevity.
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/mdm-service ./cmd/mdm-service

# Run
FROM gcr.io/distroless/base-debian12
COPY --from=builder /out/mdm-service /mdm-service
COPY configs /configs
COPY migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/mdm-service"]
