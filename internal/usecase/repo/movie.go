package repo

import (
	"context"
	"fmt"

	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/model"
	"github.com/movie-app/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MovieRepo struct {
	db     *gorm.DB
	logger *logger.Logger
	cfg    *config.Config
}

func NewMovieRepo(db *gorm.DB, cfg *config.Config, logger *logger.Logger) *MovieRepo {
	return &MovieRepo{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *MovieRepo) Create(ctx context.Context, req model.Movie) (model.Movie, error) {
	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Create(&req).Error; err != nil {
		tx.Rollback()
		r.logger.Error(fmt.Sprintf("failed to create movie: %v", err))
		return model.Movie{}, err
	}

	seen := make(map[int]struct{})
	for _, actor := range req.Cast {
		if _, ok := seen[actor.ID]; ok {
			continue
		}
		seen[actor.ID] = struct{}{}

		var existing model.Actor
		if err := tx.First(&existing, actor.ID).Error; err != nil {
			tx.Rollback()
			return model.Movie{}, fmt.Errorf("actor with ID %d not found", actor.ID)
		}

		movieActor := model.MovieActor{
			MovieID: req.ID,
			ActorID: actor.ID,
		}

		if err := tx.Table("movie_actors").
			Clauses(clause.OnConflict{DoNothing: true}).
			Create(&movieActor).Error; err != nil {
			tx.Rollback()
			return model.Movie{}, fmt.Errorf("failed to link actor ID %d to movie ID %d", actor.ID, req.ID)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return model.Movie{}, err
	}

	return req, nil
}
