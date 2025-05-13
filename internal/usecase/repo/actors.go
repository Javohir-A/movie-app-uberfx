package repo

import (
	"context"

	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/model"
	"github.com/movie-app/pkg/logger"
	"gorm.io/gorm"
)

type ActorRepo struct {
	db     *gorm.DB
	logger *logger.Logger
	cfg    *config.Config
}

func NewActorRepo(db *gorm.DB, cfg *config.Config, logger *logger.Logger) *ActorRepo {
	return &ActorRepo{
		db:     db,
		cfg:    cfg,
		logger: logger,
	}
}

func (r *ActorRepo) Create(ctx context.Context, actor model.Actor) (model.Actor, error) {
	if err := r.db.WithContext(ctx).Create(&actor).Error; err != nil {
		return model.Actor{}, err
	}
	return actor, nil
}

func (r *ActorRepo) GetByID(ctx context.Context, id uint) (model.Actor, error) {
	var actor model.Actor
	if err := r.db.WithContext(ctx).First(&actor, id).Error; err != nil {
		return model.Actor{}, err
	}
	return actor, nil
}

func (r *ActorRepo) Update(ctx context.Context, actor model.Actor) (model.Actor, error) {
	if err := r.db.WithContext(ctx).Save(&actor).Error; err != nil {
		return model.Actor{}, err
	}
	return actor, nil
}

func (r *ActorRepo) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&model.Actor{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *ActorRepo) GetList(ctx context.Context, req model.GetListFilter) (model.ActorList, error) {
	var (
		actors []model.Actor
		total  int64
	)

	tx := r.db.WithContext(ctx).Model(&model.Actor{})

	for _, filter := range req.Filters {
		query := filter.Column
		switch filter.Type {
		case "eq":
			tx = tx.Where(query+" = ?", filter.Value)
		case "ne":
			tx = tx.Where(query+" <> ?", filter.Value)
		case "gt":
			tx = tx.Where(query+" > ?", filter.Value)
		case "gte":
			tx = tx.Where(query+" >= ?", filter.Value)
		case "lt":
			tx = tx.Where(query+" < ?", filter.Value)
		case "lte":
			tx = tx.Where(query+" <= ?", filter.Value)
		case "search":
			tx = tx.Where(query+" ILIKE ?", "%"+filter.Value+"%")
		default:
			continue
		}
	}

	if err := tx.Count(&total).Error; err != nil {
		return model.ActorList{}, err
	}

	for _, order := range req.OrderBy {
		orderStr := order.Column
		if order.Order == "desc" {
			orderStr += " desc"
		} else {
			orderStr += " asc"
		}
		tx = tx.Order(orderStr)
	}

	offset := (req.Page - 1) * req.Limit
	if offset < 0 {
		offset = 0
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	if err := tx.Offset(offset).Limit(req.Limit).Find(&actors).Error; err != nil {
		return model.ActorList{}, err
	}

	return model.ActorList{
		Actors: actors,
		Total:  total,
	}, nil
}
