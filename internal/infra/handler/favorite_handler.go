package handler

import (
	"bms-go/internal/model/dto"
	"bms-go/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(s *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: s}
}

func (h *FavoriteHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/favorites")
	group.GET("", h.GetFavorites)
	group.POST("", h.AddFavorite)
}

// GetFavorites godoc
// @Summary Get all favorites
// @Description Get list of user's favorite books
// @Tags Favorites
// @Produce json
// @Success 200 {array} dto.FavoriteResponse
// @Failure 500 {object} map[string]string
// @Router /favorites [get]
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := uint(1)
	favs, err := h.service.GetFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, favs)
}

// AddFavorite godoc
// @Summary Add a favorite
// @Description Add a book to user's favorites
// @Tags Favorites
// @Accept json
// @Produce json
// @Param favorite body dto.FavoriteRequest true "Favorite request"
// @Success 201 {object} dto.FavoriteResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /favorites [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	var req dto.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := uint(1)
	resp, err := h.service.AddFavorite(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}
