package kafka

import (
	"context"
	"fmt"
	"grpc_sso/internal/config"
	"grpc_sso/internal/domain/models"
	"grpc_sso/internal/lib/logger/sl"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type Broker struct {
	log     *slog.Logger
	conn    *kafka.Conn
	options options
	Reader  *kafka.Reader
	Writer  *kafka.Writer
	ctx     context.Context
}

type options struct {
	topic             string
	address           string
	groupID           string
	network           string
	numPartitions     int
	replicationFactor int
}

func New(ctx context.Context, log *slog.Logger, cfg *config.Broker) *Broker {
	const op = "kafka.New"

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{cfg.Address},
		Topic:   cfg.TopicName,
		GroupID: cfg.GroupID,
	})

	options := options{
		network:           cfg.Network,
		address:           cfg.Address,
		topic:             cfg.TopicName,
		numPartitions:     cfg.Partitions,
		replicationFactor: cfg.Replications,
		groupID:           cfg.GroupID,
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{cfg.Address},
		Topic:   cfg.TopicName,
	})

	return &Broker{
		ctx:     ctx,
		log:     log,
		options: options,
		Reader:  r,
		Writer:  w,
	}
}

func (b *Broker) MustRun() {
	if err := b.Run(); err != nil {
		panic(err)
	}
}

func (b *Broker) Run() error {
	const op = "kafka.Run"

	log := b.log.With(
		slog.String("op", op),
		slog.String("addr", b.options.address),
		slog.String("topic name", b.options.topic),
	)

	conn, err := kafka.Dial(b.options.network, b.options.address)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             b.options.topic,
			NumPartitions:     b.options.numPartitions,
			ReplicationFactor: b.options.replicationFactor,
		},
	}

	if err := conn.CreateTopics(topicConfigs...); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("kafka is running")

	b.conn = conn
	return nil
}

func (b *Broker) Stop() {
	const op = "kafka.Stop"

	b.log.With(
		slog.String("op", op),
		slog.String("addr", b.options.address),
		slog.String("topic name", b.options.topic),
	).Info("stopping broker")

	b.conn.Close()
}

func (b *Broker) ConsumeMessage() {
	const op = "kafka.ConsumeMessage"
	for {
		msg, err := b.Reader.ReadMessage(b.ctx)
		if err != nil {
			b.log.Error("error with read message", sl.Err(err))
			continue
		}

		var sender string
		for _, header := range msg.Headers {
			if header.Key == "sender" {
				sender = string(header.Value)
				break
			}
		}

		if sender == b.options.groupID {
			b.log.Debug("skipping own message",
				slog.String("sender", sender),
				slog.String("key", string(msg.Key)))
			continue
		}

		b.log.Info("[CATCHED MESSAGE]",
			slog.String("key", string(msg.Key)),
			slog.String("value", string(msg.Value)),
			slog.Int64("offset", msg.Offset),
		)
	}
}

func (b *Broker) ProduceMessage(msg models.MessageToBroker) error {
	const op = "kafka.ProduceMessage"

	kafkaMsg := kafka.Message{
		Key:   []byte(fmt.Sprintf("key-%s", msg.Key)),
		Value: []byte(fmt.Sprintf("Message-%s", msg.Value)),
		Headers: []kafka.Header{
			{Key: "sender", Value: []byte(b.options.groupID)},
		},
	}

	err := b.Writer.WriteMessages(b.ctx, kafkaMsg)
	if err != nil {
		b.log.Error("error with sending message", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	b.log.Info("sent message successfully")
	return nil
}
