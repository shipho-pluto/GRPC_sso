package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"grpc_sso/internal/config"
	"grpc_sso/internal/domain/models"
	"grpc_sso/internal/storage"
	"grpc_sso/internal/storage/postgres/request"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

var (
	alreadyExisisCode = pqerror.Code("23505")
)

type Storage struct {
	db   *sql.DB
	info string
}

func NewStorage(storage config.Storage) *Storage {
	var pgInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		storage.Host, storage.Port, storage.User, storage.Password, storage.DBName, storage.SSLMode)
	return &Storage{
		info: pgInfo,
	}
}

func (s *Storage) Init() error {

	db, err := sql.Open("postgres", s.info)
	if err != nil {
		return fmt.Errorf("failed with open postgres db: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed with ping postgres db: %w", err)
	}

	s.db = db
	return nil
}

func (s *Storage) Close() {
	if err := s.db.Close(); err != nil {
		panic(err)
	}
}

func (s *Storage) SaveUser(
	ctx context.Context,
	email string,
	passHash []byte,
) (uid int64, err error) {
	const op = "postgres.SaveUser"

	stmt, err := s.db.Prepare(request.SaveUser)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	var id int64
	if err := stmt.QueryRowContext(ctx, email, passHash).Scan(&id); err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == alreadyExisisCode {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "postgres.User"

	stmt, err := s.db.Prepare(request.User)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	var user models.User
	if err := stmt.QueryRowContext(ctx, email).
		Scan(&user.ID, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "postgres.IsAdmin"
	stmt, err := s.db.Prepare(request.IsAdmin)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()

	var is_admin int64
	if err := stmt.QueryRowContext(ctx, userID).Scan(&is_admin); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return true, nil
}

func (s *Storage) App(ctx context.Context, appID int32) (models.App, error) {
	const op = "postgres.App"

	stmt, err := s.db.Prepare(request.App)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
		}
	}()
	var app models.App
	if err := stmt.QueryRowContext(ctx, appID).
		Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}
