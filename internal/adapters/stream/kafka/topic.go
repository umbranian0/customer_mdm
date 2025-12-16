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
	// use first broker to issue the admin request
	if len(brokers) == 0 {
		return fmt.Errorf("no brokers configured")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	ctrlConn, err := kafka.DialContext(ctx, "tcp", controller.Host+":"+fmt.Sprint(controller.Port))
	if err != nil {
		return err
	}
	defer ctrlConn.Close()

	return ctrlConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: 1,
	})
}
