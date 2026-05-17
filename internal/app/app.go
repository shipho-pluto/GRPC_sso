package app

import (
	"context"
	grpcapp "grpc_sso/internal/app/grpc"
	storapp "grpc_sso/internal/app/storage"
	"grpc_sso/internal/clients"
	"grpc_sso/internal/config"
	"grpc_sso/internal/services/auth"
	"log/slog"
	"time"
)

type App struct {
	StorageApp *storapp.App
	GRPCApp    *grpcapp.App
	Clients    *clients.Clients
}

func New(
	log *slog.Logger,
	grpcPort int,
	storage *config.DataStore,
	broker *config.Broker,
	tokenTTL time.Duration,
) *App {
	ctx := context.Background()

	storageApp := storapp.NewApp(log, storage)

	cls := clients.NewApp(ctx, log, broker)

	auth := auth.New(log, storageApp, storageApp, storageApp, cls, tokenTTL)

	grpcApp := grpcapp.NewApp(log, grpcPort, auth)

	return &App{
		Clients:    cls,
		GRPCApp:    grpcApp,
		StorageApp: storageApp,
	}
}
