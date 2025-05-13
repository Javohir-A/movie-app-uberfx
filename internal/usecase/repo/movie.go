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
func (r *MovieRepo) GetSingle(ctx context.Context, req model.Id) (model.Movie, error) {
	var movie model.Movie

	if err := r.db.WithContext(ctx).
		First(&movie, req.ID).Error; err != nil {
		return model.Movie{}, err
	}

	var cast []model.Actor
	if err := r.db.WithContext(ctx).
		Table("actors").
		Select("actors.*").
		Joins("JOIN movie_actors ON movie_actors.actor_id = actors.id").
		Where("movie_actors.movie_id = ?", movie.ID).
		Find(&cast).Error; err != nil {
		return model.Movie{}, err
	}

	movie.Cast = cast
	return movie, nil
}

func (r *MovieRepo) UpdateField(ctx context.Context, req model.UpdateFieldRequest) (model.RowsEffected, error) {
	db := r.db.WithContext(ctx).Model(&model.Movie{})

	// Apply filters
	for _, f := range req.Filter {
		switch f.Type {
		case "eq":
			db = db.Where(f.Column+" = ?", f.Value)
		case "ne":
			db = db.Where(f.Column+" <> ?", f.Value)
		case "gt":
			db = db.Where(f.Column+" > ?", f.Value)
		case "lt":
			db = db.Where(f.Column+" < ?", f.Value)
		case "gte":
			db = db.Where(f.Column+" >= ?", f.Value)
		case "lte":
			db = db.Where(f.Column+" <= ?", f.Value)
		case "search":
			db = db.Where(f.Column+" ILIKE ?", "%"+f.Value+"%")
		}
	}

	// Build update map
	updateMap := map[string]interface{}{}
	for _, item := range req.Items {
		updateMap[item.Column] = item.Value
	}

	// Execute update
	tx := db.Updates(updateMap)
	return model.RowsEffected{RowsEffected: int(tx.RowsAffected)}, tx.Error
}
