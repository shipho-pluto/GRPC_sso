package postgres

import (
	"context"
	"errors"
	"fmt"
	"grpc_sso/internal/config"
	"grpc_sso/internal/domain/models"
	"grpc_sso/internal/storage"
	"grpc_sso/internal/storage/postgres"
	"grpc_sso/internal/storage/redis"
	"log/slog"
)

type App struct {
	log       *slog.Logger
	pgDB      *postgres.Storage
	redisCl   *redis.Cache
	pgAddr    string
	redisAddr string
}

func NewApp(log *slog.Logger, storCnf *config.DataStore) *App {
	pgCfg := storCnf.Storage
	pgDB := postgres.NewStorage(pgCfg)

	redisCfg := storCnf.Cache
	redisCl := redis.NewCache(redisCfg)

	return &App{
		log:       log,
		pgDB:      pgDB,
		redisCl:   redisCl,
		pgAddr:    pgAddr(pgCfg),
		redisAddr: redisCfg.Addr,
	}
}

func (a *App) Run() error {
	const op = "storage.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.String("postgres address", a.pgAddr),
		slog.String("redis address", a.redisAddr),
	)

	if err := a.pgDB.Init(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("postgres is running")

	if err := a.redisCl.Init(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("redis is running")

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	const op = "storage.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping storage",
			slog.String("postgres address", a.pgAddr),
			slog.String("redis address", a.redisAddr),
		)

	a.pgDB.Close()
	a.redisCl.Close()
}

func (a *App) SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error) {
	return a.pgDB.SaveUser(ctx, email, passHash)
}
func (a *App) User(ctx context.Context, email string) (models.User, error) {
	user, err := a.redisCl.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrorNotInRedis) {
			user, err = a.pgDB.User(ctx, email)
			if err != nil {
				return user, nil
			}
			if err := a.redisCl.CacheUser(ctx, user); err != nil {
				return models.User{}, err
			}
			return user, nil
		}
		return models.User{}, err
	}
	return user, nil
}
func (a *App) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	is_admin, err := a.redisCl.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotInRedis) {
			is_admin, err = a.pgDB.IsAdmin(ctx, userID)
			if err != nil {
				return is_admin, nil
			}
			if err := a.redisCl.CacheIsAdmin(ctx, userID, is_admin); err != nil {
				return false, err
			}
			return is_admin, nil
		}
		return false, err
	}
	return is_admin, nil
}
func (a *App) App(ctx context.Context, appID int32) (models.App, error) {
	app, err := a.redisCl.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrorNotInRedis) {
			app, err = a.pgDB.App(ctx, appID)
			if err != nil {
				return app, nil
			}
			if err := a.redisCl.CacheApp(ctx, app); err != nil {
				return models.App{}, err
			}
			return app, nil
		}
		return models.App{}, err
	}
	return app, nil
}

func pgAddr(cfg config.Storage) string {
	return cfg.Host + ":" + cfg.Port
}
