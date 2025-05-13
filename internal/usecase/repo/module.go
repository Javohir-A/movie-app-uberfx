package repo

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewMovieRepo),
	fx.Provide(NewActorRepo),
)
