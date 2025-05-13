package router

import (
	"context"

	"github.com/gin-gonic/gin"
	_ "github.com/movie-app/docs"
	"github.com/movie-app/internal/config"
	"go.uber.org/fx"
)

func NewRouter() *gin.Engine {
	return gin.Default()
}

func RegisterHooks(lc fx.Lifecycle, router *gin.Engine, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := router.Run(":" + cfg.Port); err != nil {
					panic("Failed to start server: " + err.Error())
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// Optional: graceful shutdown
			return nil
		},
	})
}
