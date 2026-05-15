package redis

import (
	"context"
	"errors"
	"fmt"
	"grpc_sso/internal/config"
	"grpc_sso/internal/domain/models"
	"grpc_sso/internal/storage"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var (
	expirationTime = time.Minute
)

type Cache struct {
	cl  *redis.Client
	opt *redis.Options
}

func NewCache(cfg config.Cache) *Cache {
	var opt = &redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	}

	return &Cache{
		opt: opt,
	}
}

func (c *Cache) Init() error {

	cl := redis.NewClient(c.opt)

	if err := cl.Ping().Err(); err != nil {
		return fmt.Errorf("failed with ping redis: %w", err)
	}

	c.cl = cl
	return nil
}

func (c *Cache) Close() {
	if err := c.cl.Close(); err != nil {
		panic(err)
	}
}

func (c *Cache) CacheUser(ctx context.Context, user models.User) error {
	const op = "redis.CacheUser"
	err := c.cl.Set(user.Email, user, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) CacheApp(ctx context.Context, app models.App) error {
	const op = "redis.CacheApp"
	err := c.cl.Set(i32ToS(app.ID), app, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) CacheIsAdmin(ctx context.Context, uid int64, is_admin bool) error {
	const op = "redis.CacheIsAdmin"

	err := c.cl.Set(i64ToS(uid), is_admin, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Cache) User(ctx context.Context, email string) (models.User, error) {
	const op = "redis.User"

	var user models.User
	if err := c.cl.Get(email).Scan(&user); err != nil {
		if errors.Is(err, redis.Nil) {
			return models.User{}, storage.ErrorNotInRedis
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (c *Cache) App(ctx context.Context, app_id int32) (models.App, error) {
	const op = "redis.App"

	var app models.App
	if err := c.cl.Get(i32ToS(app_id)).Scan(&app); err != nil {
		if errors.Is(err, redis.Nil) {
			return models.App{}, storage.ErrorNotInRedis
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (c *Cache) IsAdmin(ctx context.Context, uid int64) (bool, error) {
	const op = "redis.IsAdmin"

	var is_admin bool
	if err := c.cl.Get(i64ToS(uid)).Scan(&is_admin); err != nil {
		if errors.Is(err, redis.Nil) {
			return false, storage.ErrorNotInRedis
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return is_admin, nil
}

func i64ToS(x int64) string {
	return strconv.Itoa(int(x))
}

func i32ToS(x int32) string {
	return strconv.Itoa(int(x))
}
