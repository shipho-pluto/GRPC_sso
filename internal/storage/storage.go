package storage

import "errors"

var (
	ErrUserExists   = errors.New("user already exisis")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
	ErrorNotInRedis = errors.New("not in redis")
)
