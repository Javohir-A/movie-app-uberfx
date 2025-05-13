package logger

import (
	"github.com/movie-app/internal/config"
	"go.uber.org/fx"
)

func SetupLogger(cfg *config.Config) *Logger {
	return New(cfg.LogLevel)
}

var Module = fx.Option(
	fx.Provide(SetupLogger),
)
