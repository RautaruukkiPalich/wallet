package kafka

import (
	"context"
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

func NewConsumer(cfg ConsumerConfig) *Consumer {

	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers: []string{cfg.Addr},
			Topic:   cfg.Topic,
			GroupID: cfg.GroupID,
		})

	return &Consumer{r: r}
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
