package web

import (
	"log/slog"
	"os"
)

func NewLogger(name string) *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}).WithAttrs([]slog.Attr{{Key: "service", Value: slog.StringValue(name)}}))
}
