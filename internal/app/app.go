package app

import (
	grpcapp "grpc_sso/internal/app/grpc"
	storapp "grpc_sso/internal/app/storage"
	"grpc_sso/internal/config"
	"grpc_sso/internal/services/auth"
	"log/slog"
	"time"
)

type App struct {
	StorageApp *storapp.App
	GRPCApp    *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storage *config.DataStore,
	tokenTTL time.Duration,
) *App {

	storageApp := storapp.NewApp(log, storage)

	auth := auth.New(log, storageApp, storageApp, storageApp, tokenTTL)

	grpcApp := grpcapp.NewApp(log, grpcPort, auth)

	return &App{
		GRPCApp:    grpcApp,
		StorageApp: storageApp,
	}
}
