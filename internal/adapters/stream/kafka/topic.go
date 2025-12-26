package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// ensureTopic creates the topic if it does not exist. Safe to call repeatedly.
func ensureTopic(brokers []string, topic string, partitions int) error {
	if partitions <= 0 {
		partitions = 1
	}
	if len(brokers) == 0 {
		return fmt.Errorf("no brokers configured")
	}

	var lastErr error
	for _, broker := range brokers {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		conn, err := kafka.DialContext(ctx, "tcp", broker)
		if err != nil {
			lastErr = err
			cancel()
			continue
		}

		controller, err := conn.Controller()
		if err != nil {
			lastErr = err
			conn.Close()
			cancel()
			continue
		}

		ctrlConn, err := kafka.DialContext(ctx, "tcp", controller.Host+":"+fmt.Sprint(controller.Port))
		if err != nil {
			lastErr = err
			conn.Close()
			cancel()
			continue
		}

		err = ctrlConn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1,
		})
		ctrlConn.Close()
		conn.Close()
		cancel()
		if err == nil {
			return nil
		}
		lastErr = err
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("failed to ensure topic %s", topic)
	}
	return lastErr
}
