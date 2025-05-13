package handler

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/model"
	"github.com/movie-app/internal/usecase"
	"github.com/movie-app/pkg/logger"
	"github.com/spf13/cast"
)

type MovieHandler struct {
	usecase *usecase.UseCase
	cfg     *config.Config
	logger  *logger.Logger
}

func NewMovieHandler(usecase *usecase.UseCase, cfg *config.Config, logger *logger.Logger) *MovieHandler {
	return &MovieHandler{
		usecase: usecase,
		logger:  logger,
		cfg:     cfg,
	}
}

func (h *MovieHandler) RegisterRoutes(r *gin.Engine) {
	movieHandler := r.Group("/v1/movies")
	{
		movieHandler.POST("", h.Create)
		movieHandler.GET("/:id", h.GetByID)
		movieHandler.PUT("/:id", h.Update)
		movieHandler.GET("", h.GetAll)
		movieHandler.PUT("/field")
		movieHandler.DELETE("/:id", h.Delete)
	}
}

// @Summary Create a new movie
// @Description Create a new movie with the given JSON payload
// @Tags movies
// @Accept json
// @Produce json
// @Param movie body model.CreateMovieRequest true "Movie data"
// @Success 201 {object} model.Movie
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/movies [post]
func (h *MovieHandler) Create(c *gin.Context) {
	var (
		movie model.CreateMovieRequest
		cast  []model.Actor
	)

	if err := c.ShouldBindJSON(&movie); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]string, len(ve))
			for i, fe := range ve {
				out[i] = fmt.Sprintf("Field '%s' is %s", fe.Field(), fe.Tag())
			}
			c.JSON(400, gin.H{"error": out})
		} else {
			h.logger.Error(fmt.Sprintf("Failed to bind movie: %v", err))
			c.JSON(400, gin.H{"error": "Invalid movie data"})
		}
		return
	}

	for _, castInput := range movie.Casts {
		cast = append(cast, model.Actor{ID: castInput.Id})
	}

	createdMovie, err := h.usecase.MovieRepo.Create(c.Request.Context(), model.Movie{
		Title:    movie.Title,
		Director: movie.Director,
		Plot:     movie.Plot,
		Year:     movie.Year,
		Cast:     cast,
	})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to create movie: %v", err))
		c.JSON(500, gin.H{"error": "Failed to create movie"})
		return
	}

	c.JSON(201, createdMovie)
}

// @Summary Get movie by ID
// @Description Fetch a single movie by its ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} model.Movie
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /v1/movies/{id} [get]
func (h *MovieHandler) GetByID(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSON(400, gin.H{"error": "id must be provided"})
		return
	}
	movie, err := h.usecase.MovieRepo.GetSingle(c.Request.Context(), model.Id{ID: id})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to fetch movie by ID: %v", err))
		c.JSON(404, gin.H{"error": "Movie not found"})
		return
	}
	c.JSON(200, movie)
}

// @Summary Update movie
// @Description Update a movie by its ID
// @Tags movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body model.UpdateMovieRequest true "Updated movie"
// @Success 200 {object} model.Movie
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/movies/{id} [put]
func (h *MovieHandler) Update(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSON(400, gin.H{"error": "id must be provided"})
		return
	}

	var req model.UpdateMovieRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error(fmt.Sprintf("Failed to bind update movie data: %v", err))
		c.JSON(400, gin.H{"error": "Invalid movie data"})
		return
	}

	var cast []model.Actor
	for _, actor := range req.Casts {
		cast = append(cast, model.Actor{ID: actor.Id})
	}

	movie := model.Movie{
		ID:       id,
		Title:    req.Title,
		Director: req.Director,
		Plot:     req.Plot,
		Year:     req.Year,
		Cast:     cast,
	}

	updatedMovie, err := h.usecase.MovieRepo.Update(c.Request.Context(), movie)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to update movie: %v", err))
		c.JSON(500, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(200, updatedMovie)
}

// @Summary Delete movie
// @Description Delete a movie by its ID
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} map[string]any{}
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/movies/{id} [delete]
func (h *MovieHandler) Delete(c *gin.Context) {
	id := cast.ToInt(c.Param("id"))
	if id == 0 {
		c.JSON(400, gin.H{"error": "id must be provided"})
		return
	}
	err := h.usecase.MovieRepo.Delete(c.Request.Context(), model.Id{ID: id})
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to delete movie: %v", err))
		c.JSON(500, gin.H{"error": "Failed to delete movie"})
		return
	}
	c.JSON(200, gin.H{"message": "Movie deleted successfully"})
}

// @Summary Get all movies
// @Description Get a paginated list of movies with optional filters and ordering
// @Tags movies
// @Produce json
// @Param page query int false "Page number (default is 1)"
// @Param limit query int false "Items per page (default is 10)"
// @Param title query string false "Search by movie title"
// @Param director query string false "Search by director name"
// @Param year query string false "Search by release year"
// @Param order_by query string false "Field to order by (e.g. year)"
// @Param sort query string false "Sort direction: asc or desc (default asc)"
// @Success 200 {object} model.MovieList
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/movies [get]
func (h *MovieHandler) GetAll(c *gin.Context) {
	var req model.GetListFilter

	req.Page = parseInt(c.DefaultQuery("page", "1"), 1)
	req.Limit = parseInt(c.DefaultQuery("limit", "10"), 10)

	for key, values := range c.Request.URL.Query() {
		switch key {
		case "title", "director", "year":
			req.Filters = append(req.Filters, model.Filter{
				Column: key,
				Type:   "search",
				Value:  values[0],
			})
		}
	}

	orderBy := c.Query("order_by")
	sort := c.DefaultQuery("sort", "asc")
	if orderBy != "" {
		req.OrderBy = append(req.OrderBy, model.OrderBy{
			Column: orderBy,
			Order:  sort,
		})
	}

	// Call repo
	movies, err := h.usecase.MovieRepo.GetList(c.Request.Context(), req)
	if err != nil {
		h.logger.Error(fmt.Sprintf("Failed to fetch movies: %v", err))
		c.JSON(500, gin.H{"error": "Failed to fetch movies"})
		return
	}

	c.JSON(200, movies)
}

// func (h *MovieHandler) UpdateField(c *gin.Context) {
// 	// Expecting JSON like: { "id": "movie_id", "field": "title", "value": "New Title" }
// 	var req struct {
// 		ID    string `json:"id"`
// 		Field string `json:"field"`
// 		Value string `json:"value"`
// 	}

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		h.logger.Error(fmt.Sprintf("Invalid update field payload: %v", err))
// 		c.JSON(400, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	err := h.usecase.MovieRepo.UpdateField(c.Request.Context(), req.ID, req.Field, req.Value)
// 	if err != nil {
// 		h.logger.Errorf("Failed to update field: %v", err)
// 		c.JSON(500, gin.H{"error": "Failed to update movie field"})
// 		return
// 	}
// 	c.JSON(200, gin.H{"message": "Field updated successfully"})
// }
