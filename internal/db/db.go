package db

import (
	"fmt"
	"log"

	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/movie-app/internal/config"
)

func NewGormDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	log.Println(dsn)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}

	return db, nil
}

var Module = fx.Option(
	fx.Provide(NewGormDatabase),
)
