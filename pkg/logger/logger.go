package logger

import (
	"github.com/lmittmann/tint"
	"log/slog"
	"os"
)

var Log = NewLogger()

func NewLogger() *slog.Logger {
	opts := &tint.Options{
		Level:      slog.LevelDebug,
		TimeFormat: "2006-01-02 15:04:05.000",
	}
	handler := tint.NewHandler(os.Stdout, opts)
	logger := slog.New(handler)

	return logger
}
