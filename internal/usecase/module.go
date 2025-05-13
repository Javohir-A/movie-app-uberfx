package usecase

import (
	"github.com/movie-app/internal/usecase/repo"
	"go.uber.org/fx"
)

func provideMovieRepoInterface(r *repo.MovieRepo) MovieRepoI {
	return r
}
func provideActorRepoInterface(r *repo.ActorRepo) ActorRepoI {
	return r
}

var Module = fx.Options(
	repo.Module,
	fx.Provide(
		provideMovieRepoInterface,
		provideActorRepoInterface,
		NewUseCase,
	),
)
