package handler

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewMovieHandler),
	fx.Provide(NewActorHandler),
)
