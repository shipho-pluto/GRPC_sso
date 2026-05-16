package main

import (
	"grpc_sso/internal/app"
	"grpc_sso/internal/config"
	"grpc_sso/internal/lib/logger/setup"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.MustLoad()

	log := setup.SetupLogger(cfg.Env)

	log.Info("starting application...", slog.String("env", cfg.Env))
	log.Debug("start debug logger")

	application := app.New(log, cfg.GRPC.Port, &cfg.DataStore, &cfg.Clients.Broker, cfg.TokenTTL)

	go application.GRPCApp.MustRun()
	go application.StorageApp.MustRun()
	go application.Clients.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop
	log.Debug("stopping aplication...", slog.String("signal", sign.String()))

	application.GRPCApp.Stop()
	application.StorageApp.Stop()
	application.Clients.Stop()
	log.Info("application stopped")

	//TODO: kafka between sso & url
	//TODO: tests
}
