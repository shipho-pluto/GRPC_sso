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
}

type options struct {
	topic             string
	address           string
	network           string
	numPartitions     int
	replicationFactor int
}

func New(log *slog.Logger, cfg *config.Broker) *Broker {
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
	}

	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{cfg.Address},
		Topic:   cfg.TopicName,
	})

	return &Broker{
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
	)

	conn, err := kafka.Dial(b.options.network, b.options.address)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("kafla is created")

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

	b.log.With(slog.String("op", op)).
		Info("stopping broker")

	b.conn.Close()
}

func (b *Broker) ConsumeMessage(ctx context.Context) {
	const op = "kafka.ConsumeMessage"
	for {
		msg, err := b.Reader.ReadMessage(ctx)
		if err != nil {
			b.log.Error("error with read message", sl.Err(err))
		}

		b.log.Info("[CATCHED MESSAGE]",
			slog.String("key", string(msg.Key)),
			slog.String("value", string(msg.Value)),
			slog.Int64("offset", msg.Offset),
		)
	}
}

func (b *Broker) ProduceMessage(ctx context.Context, msg models.MessageToBroker) error {
	const op = "kafka.ProduceMessage"

	kafkaMsg := kafka.Message{
		Key:   []byte(fmt.Sprintf("key-%s", msg.Key)),
		Value: []byte(fmt.Sprintf("Mwssage-%s", msg.Value)),
	}

	err := b.Writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		b.log.Error("error with sending message", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	b.log.Info("sent message successfully")
	return nil
}
