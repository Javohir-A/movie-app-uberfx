package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/movie-app/docs"
	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/handler"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// NewRouter -.
// Swagger spec:
// @title       Movie APIs
// @description This is a movie CRUD APIs
// @version     2.0
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func SetupRoutes(
	router *gin.Engine,
	movieHandler *handler.MovieHandler,
	actorHandler *handler.ActorHandler,
) {

	// Swagger
	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Prometheus metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	movieHandler.RegisterRoutes(router)
	actorHandler.RegisterRoutes(router)
}

var Module = fx.Options(
	fx.Provide(NewRouter),
	fx.Invoke(RegisterHooks),
	fx.Invoke(SetupRoutes),
)
