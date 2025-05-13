package app

import (
	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/db"
	"github.com/movie-app/internal/handler"
	"github.com/movie-app/internal/router"
	"github.com/movie-app/internal/usecase"
	"github.com/movie-app/pkg/logger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	logger.Module,
	db.Module,
	usecase.Module,
	handler.Module,
	router.Module,
)
