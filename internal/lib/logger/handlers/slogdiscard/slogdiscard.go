package slogdiscard

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func NewDiscardHandler() *DiscardHandler {
	return &DiscardHandler{}
}

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

func (d *DiscardHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (d *DiscardHandler) Handle(context.Context, slog.Record) error {
	return nil
}

func (d *DiscardHandler) WithAttrs([]slog.Attr) slog.Handler {
	return d
}

func (d *DiscardHandler) WithGroup(string) slog.Handler {
	return d
}
