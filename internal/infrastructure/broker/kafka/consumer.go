package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type ConsumerConfig struct {
	Addr    string
	Topic   string
	GroupID string
}

type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(cfg ConsumerConfig) (*Consumer, error) {

	conn, err := kafka.DialContext(context.Background(), "tcp", cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Kafka: %w", err)
	}
	defer conn.Close()

	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{cfg.Addr},
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
		})

	return &Consumer{r: r}, nil
}

func (c *Consumer) Consume(ctx context.Context) ([]byte, error) {
	msg, err := c.r.ReadMessage(ctx)
	if err != nil {
		return nil, err
	}
	return msg.Value, nil
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
