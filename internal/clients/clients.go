package clients

import (
	"context"
	"grpc_sso/internal/clients/kafka"
	"grpc_sso/internal/config"
	"log/slog"
)

type Clients struct {
	*kafka.Broker
	log *slog.Logger
}

func NewApp(ctx context.Context, log *slog.Logger, cfg *config.Broker) *Clients {
	broker := kafka.New(ctx, log, cfg)
	return &Clients{
		log:    log,
		Broker: broker,
	}
}

func (c *Clients) MustRun() {
	c.Broker.MustRun()
}

func (c *Clients) Stop() {
	c.Broker.Stop()
}
