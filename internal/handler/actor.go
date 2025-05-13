package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/movie-app/internal/config"
	"github.com/movie-app/internal/model"
	"github.com/movie-app/internal/usecase"
	"github.com/movie-app/pkg/logger"
)

type ActorHandler struct {
	usecase *usecase.UseCase
	cfg     *config.Config
	logger  *logger.Logger
}

func NewActorHandler(usecase *usecase.UseCase, cfg *config.Config, logger *logger.Logger) *ActorHandler {
	return &ActorHandler{
		usecase: usecase,
		logger:  logger,
		cfg:     cfg,
	}
}

func (h *ActorHandler) RegisterRoutes(r *gin.Engine) {
	actorHandler := r.Group("/v1/actors")
	{
		actorHandler.POST("", h.Create)
		actorHandler.GET("/:id", h.GetByID)
		actorHandler.PUT("/:id", h.Update)
		actorHandler.GET("", h.GetList)
		actorHandler.DELETE("/:id", h.Delete)
	}
}

// Create godoc
// @Summary Create a new actor
// @Description Creates an actor and returns the created object
// @Tags actors
// @Accept json
// @Produce json
// @Param actor body model.Actor true "Actor data"
// @Success 201 {object} model.Actor
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/actors [post]
func (h *ActorHandler) Create(c *gin.Context) {
	var req model.Actor
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid actor payload: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid payload", Code: "BAD_REQUEST"})
		return
	}

	res, err := h.usecase.ActorRepo.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to create actor: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to create actor", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// GetByID godoc
// @Summary Get actor by ID
// @Description Retrieves an actor by its ID
// @Tags actors
// @Produce json
// @Param id path int true "Actor ID"
// @Success 200 {object} model.Actor
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /v1/actors/{id} [get]
func (h *ActorHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("invalid actor id: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid actor ID", Code: "BAD_REQUEST"})
		return
	}

	actor, err := h.usecase.ActorRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		h.logger.Error("actor not found: %v", err)
		c.JSON(http.StatusNotFound, model.ErrorResponse{Message: "Actor not found", Code: "NOT_FOUND"})
		return
	}

	c.JSON(http.StatusOK, actor)
}

// Update godoc
// @Summary Update an actor
// @Description Updates actor information by ID
// @Tags actors
// @Accept json
// @Produce json
// @Param id path int true "Actor ID"
// @Param actor body model.Actor true "Actor data"
// @Success 200 {object} model.Actor
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/actors/{id} [put]
func (h *ActorHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("invalid actor id: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid actor ID", Code: "BAD_REQUEST"})
		return
	}

	var req model.Actor
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("invalid update payload: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid payload", Code: "BAD_REQUEST"})
		return
	}
	req.ID = id

	updated, err := h.usecase.ActorRepo.Update(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to update actor: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to update actor", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete godoc
// @Summary Delete an actor
// @Description Deletes an actor by ID
// @Tags actors
// @Param id path int true "Actor ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/actors/{id} [delete]
func (h *ActorHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("invalid actor id: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid actor ID", Code: "BAD_REQUEST"})
		return
	}

	if err := h.usecase.ActorRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		h.logger.Error("failed to delete actor: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to delete actor", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, model.SuccessResponse{Message: "Actor deleted successfully"})
}

// GetList godoc
// @Summary Get a list of actors
// @Description Retrieves a paginated list of actors
// @Tags actors
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Page size"
// @Success 200 {object} model.ActorList
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /v1/actors [get]
func (h *ActorHandler) GetList(c *gin.Context) {
	var req model.GetListFilter
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Error("invalid query params: %v", err)
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Invalid query params", Code: "BAD_REQUEST"})
		return
	}

	list, err := h.usecase.ActorRepo.GetList(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to get actor list: %v", err)
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: "Failed to fetch actor list", Code: "INTERNAL_ERROR"})
		return
	}

	c.JSON(http.StatusOK, list)
}
