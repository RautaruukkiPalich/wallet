package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type ProducerConfig struct {
	Addr  string
	Topic string
}

type Producer struct {
	pr *kafka.Writer
}

func NewProducer(cfg ProducerConfig) (*Producer, error) {

	// hint for fix: panic: [3] Unknown Topic Or Partition: the request is for a topic or partition that does not exist on this broker
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		cfg.Addr,
		cfg.Topic,
		0,
	)
	if err != nil {
		return nil, err
	}
	// close the connection because we won't be using it
	_ = conn.Close()

	p := &Producer{
		pr: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Addr),
			Topic:        cfg.Topic,
			BatchBytes:   0,
			BatchSize:    1, //
			BatchTimeout: 0,
		},
	}
	return p, nil
}

func (p *Producer) Publish(ctx context.Context, data []byte) error {
	msg := kafka.Message{
		Value: data,
	}

	return p.pr.WriteMessages(ctx, msg)
}

func (p *Producer) Close() error {
	return p.pr.Close()
}
