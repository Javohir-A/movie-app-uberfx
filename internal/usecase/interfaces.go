package usecase

import (
	"context"

	"github.com/movie-app/internal/model"
)

type (
	MovieRepoI interface {
		Create(ctx context.Context, req model.Movie) (model.Movie, error)
		GetSingle(ctx context.Context, req model.Id) (model.Movie, error)
		UpdateField(ctx context.Context, req model.UpdateFieldRequest) (model.RowsEffected, error)
		Update(ctx context.Context, req model.Movie) (model.Movie, error)
		Delete(ctx context.Context, req model.Id) error
		GetList(ctx context.Context, req model.GetListFilter) (model.MovieList, error)
	}

	ActorRepoI interface {
		Create(ctx context.Context, actor model.Actor) (model.Actor, error)
		GetByID(ctx context.Context, id uint) (model.Actor, error)
		Update(ctx context.Context, actor model.Actor) (model.Actor, error)
		Delete(ctx context.Context, id uint) error
		GetList(ctx context.Context, req model.GetListFilter) (model.ActorList, error)
	}
)
