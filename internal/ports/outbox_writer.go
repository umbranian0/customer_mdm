package ports

type OutboxWriter interface {
    Write(tx Tx, topic string, key, value []byte, headers map[string]string) error
}
