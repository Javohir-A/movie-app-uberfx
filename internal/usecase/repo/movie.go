package repo

import (
	"context"
	"fmt"

	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/model"
	"github.com/movie-app/pkg/logger"
	"go.uber.org/zap"
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
func (r *MovieRepo) Update(ctx context.Context, req model.Movie) (model.Movie, error) {
	tx := r.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&model.Movie{}).
		Where("id = ?", req.ID).
		Updates(map[string]any{
			"title":      req.Title,
			"director":   req.Director,
			"year":       req.Year,
			"plot":       req.Plot,
			"updated_at": gorm.Expr("NOW()"),
		}).Error; err != nil {
		tx.Rollback()
		return model.Movie{}, fmt.Errorf("failed to update movie fields: %w", err)
	}

	// clear old cast (relations)
	if err := tx.Where("movie_id = ?", req.ID).Delete(&model.MovieActor{}).Error; err != nil {
		tx.Rollback()
		return model.Movie{}, fmt.Errorf("failed to clear old cast: %w", err)
	}

	// Validate adn recreate actor links
	seen := make(map[int]struct{})
	for _, actor := range req.Cast {
		if _, exists := seen[actor.ID]; exists {
			continue // avoiding duplicates
		}
		seen[actor.ID] = struct{}{}

		//  actor exists?
		var exists int64
		if err := tx.Model(&model.Actor{}).
			Where("id = ?", actor.ID).
			Count(&exists).Error; err != nil {
			tx.Rollback()
			return model.Movie{}, fmt.Errorf("failed to validate actor ID %d: %w", actor.ID, err)
		}
		if exists == 0 {
			tx.Rollback()
			return model.Movie{}, fmt.Errorf("actor with ID %d not found", actor.ID)
		}

		// Create movie-actor relation
		link := model.MovieActor{MovieID: req.ID, ActorID: actor.ID}
		if err := tx.Create(&link).Error; err != nil {
			tx.Rollback()
			return model.Movie{}, fmt.Errorf("failed to link actor ID %d: %w", actor.ID, err)
		}
	}

	// fetch updated movie with preloaded cast
	var updated model.Movie
	if err := tx.Preload("Cast").First(&updated, req.ID).Error; err != nil {
		tx.Rollback()
		return model.Movie{}, fmt.Errorf("failed to fetch updated movie: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return model.Movie{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updated, nil
}

func (r *MovieRepo) GetList(ctx context.Context, req model.GetListFilter) (model.MovieList, error) {
	var movies []model.Movie
	query := r.db.WithContext(ctx).Model(&model.Movie{})

	// Filters
	for _, f := range req.Filters {
		switch f.Type {
		case "eq":
			query = query.Where(f.Column+" = ?", f.Value)
		case "search":
			query = query.Where(f.Column+" ILIKE ?", "%"+f.Value+"%")
		case "gt":
			query = query.Where(f.Column+" > ?", f.Value)
		case "lt":
			query = query.Where(f.Column+" < ?", f.Value)
		case "gte":
			query = query.Where(f.Column+" >= ?", f.Value)
		case "lte":
			query = query.Where(f.Column+" <= ?", f.Value)
		}
	}

	// Ordering
	for _, o := range req.OrderBy {
		query = query.Order(fmt.Sprintf("%s %s", o.Column, o.Order))
	}

	// Count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return model.MovieList{}, err
	}

	// Pagination
	offset := (req.Page - 1) * req.Limit
	if offset < 0 {
		offset = 0
	}

	if err := query.Limit(req.Limit).Offset(offset).Find(&movies).Error; err != nil {
		return model.MovieList{}, err
	}

	// Attach cast manually
	for i := range movies {
		var cast []model.Actor
		err := r.db.Table("actors").
			Select("actors.*").
			Joins("JOIN movie_actors ON movie_actors.actor_id = actors.id").
			Where("movie_actors.movie_id = ?", movies[i].ID).
			Find(&cast).Error
		if err != nil {
			return model.MovieList{}, err
		}
		movies[i].Cast = cast
	}

	return model.MovieList{
		Movies: movies,
		Count:  int(total),
	}, nil
}
func (r *MovieRepo) Delete(ctx context.Context, req model.Id) error {
	tx := r.db.WithContext(ctx).Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("movie_id = ?", req.ID).Delete(&model.MovieActor{}).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to delete movie_actors relations", zap.Int("movie_id", req.ID), zap.Error(err))
		return fmt.Errorf("failed to delete movie_actors relations for movie ID %d: %w", req.ID, err)
	}

	if err := tx.Delete(&model.Movie{}, req.ID).Error; err != nil {
		tx.Rollback()
		r.logger.Error("failed to delete movie", zap.Int("movie_id", req.ID), zap.Error(err))
		return fmt.Errorf("failed to delete movie ID %d: %w", req.ID, err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.logger.Error("transaction commit failed", zap.Int("movie_id", req.ID), zap.Error(err))
		return fmt.Errorf("transaction commit failed for movie ID %d: %w", req.ID, err)
	}

	// Successful deletion
	return nil
}
