package handler

import (
	"bms-go/internal/model/dto"
	"bms-go/internal/service"
	"net/http"
	"strconv"

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
	group.GET("/:id", h.GetFavoriteByID)
	group.POST("", h.AddFavorite)
	group.DELETE("/:id", h.RemoveFavorite)
}

// GetFavorites godoc
// @Summary Get all user favorites
// @Description Retrieve a list of all books marked as favorites by the current user. Each favorite includes book details such as title, author, and category. The user ID is currently hardcoded to 1 for demo purposes.
// @Tags Favorites
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse{data=[]model.SwaggerFavorite} "Favorites retrieved successfully"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /favorites [get]
// @Example {
//   "request": "GET /favorites",
//   "response": {
//     "success": true,
//     "message": "Favorites retrieved successfully",
//     "data": [
//       {
//         "id": 1,
//         "user_id": 1,
//         "book_id": 1,
//         "created_at": "2023-01-01T00:00:00Z",
//         "book": {
//           "id": 1,
//           "title": "Harry Potter and the Sorcerer's Stone",
//           "author": "J.K. Rowling",
//           "category": "Fantasy",
//           "created_at": "2023-01-01T00:00:00Z",
//           "updated_at": "2023-01-01T00:00:00Z"
//         }
//       }
//     ]
//   }
// }
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	userID := h.getUserIDFromContext(c)
	
	favs, err := h.service.GetFavorites(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve favorites",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Favorites retrieved successfully",
		Data:    favs,
	})
}

// GetFavoriteByID godoc
// @Summary Get favorite by ID
// @Description Retrieve detailed information about a specific favorite entry using its unique identifier. Returns complete favorite details including the associated book information.
// @Tags Favorites
// @Accept json
// @Produce json
// @Param id path int true "Favorite ID (must be a positive integer)"
// @Success 200 {object} dto.APIResponse{data=model.SwaggerFavorite} "Favorite retrieved successfully"
// @Failure 400 {object} dto.APIResponse "Invalid favorite ID format"
// @Failure 404 {object} dto.APIResponse "Favorite not found"
// @Router /favorites/{id} [get]
// @Example {
//   "request": "GET /favorites/1",
//   "response": {
//     "success": true,
//     "message": "Favorite retrieved successfully",
//     "data": {
//       "id": 1,
//       "user_id": 1,
//       "book_id": 1,
//       "created_at": "2023-01-01T00:00:00Z",
//       "book": {
//         "id": 1,
//         "title": "Harry Potter and the Sorcerer's Stone",
//         "author": "J.K. Rowling",
//         "category": "Fantasy",
//         "created_at": "2023-01-01T00:00:00Z",
//         "updated_at": "2023-01-01T00:00:00Z"
//       }
//     }
//   }
// }
func (h *FavoriteHandler) GetFavoriteByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid favorite ID",
			Error:   "Favorite ID must be a positive integer",
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	fav, err := h.service.GetFavoriteByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Favorite not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Favorite retrieved successfully",
		Data:    fav,
	})
}

// AddFavorite godoc
// @Summary Add book to favorites
// @Description Add a book to the current user's favorites list. Validates that the book exists and checks for duplicates. The user ID is currently hardcoded to 1 for demo purposes. Returns the complete favorite details including book information.
// @Tags Favorites
// @Accept json
// @Produce json
// @Param favorite body dto.FavoriteRequest true "Favorite request" required(true)
// @Success 201 {object} dto.APIResponse{data=dto.FavoriteResponse} "Favorite added successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request body or validation failed"
// @Failure 404 {object} dto.APIResponse "Book not found"
// @Failure 409 {object} dto.APIResponse "Book already in favorites"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /favorites [post]
// @Example {
//   "request": {
//     "book_id": 1
//   },
//   "response": {
//     "success": true,
//     "message": "Favorite added successfully",
//     "data": {
//       "id": 1,
//       "user_id": 1,
//       "book_id": 1,
//       "created_at": "2023-01-01T00:00:00Z",
//       "book": {
//         "id": 1,
//         "title": "Harry Potter and the Sorcerer's Stone",
//         "author": "J.K. Rowling",
//         "category": "Fantasy",
//         "created_at": "2023-01-01T00:00:00Z",
//         "updated_at": "2023-01-01T00:00:00Z"
//       }
//     }
//   }
// }
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	var req dto.FavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate input
	if err := h.validateFavoriteRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	resp, err := h.service.AddFavorite(userID, req)
	if err != nil {
		// Handle specific error cases
		if err.Error() == "book not found" {
			c.JSON(http.StatusNotFound, dto.APIResponse{
				Success: false,
				Message: "Book not found",
				Error:   "The specified book does not exist",
			})
			return
		}
		if err.Error() == "already in favorites" {
			c.JSON(http.StatusConflict, dto.APIResponse{
				Success: false,
				Message: "Book already in favorites",
				Error:   "This book is already in your favorites",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to add favorite",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Favorite added successfully",
		Data:    resp,
	})
}

// RemoveFavorite godoc
// @Summary Remove a favorite
// @Description Remove a book from user's favorites
// @Tags Favorites
// @Produce json
// @Param id path int true "Favorite ID"
// @Success 200 {object} dto.APIResponse
// @Failure 400 {object} dto.APIResponse
// @Failure 404 {object} dto.APIResponse
// @Failure 500 {object} dto.APIResponse
// @Router /favorites/{id} [delete]
func (h *FavoriteHandler) RemoveFavorite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid favorite ID",
			Error:   "Favorite ID must be a positive integer",
		})
		return
	}

	userID := h.getUserIDFromContext(c)
	err = h.service.RemoveFavorite(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Favorite not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Favorite removed successfully",
	})
}

// getUserIDFromContext extracts user ID from context
// For now, returns hardcoded user ID 1 as per original implementation
// In a real app, this would extract from JWT token or session
func (h *FavoriteHandler) getUserIDFromContext(c *gin.Context) uint {
	// TODO: Extract from JWT token or session in production
	return uint(1)
}

// validateFavoriteRequest validates the favorite request data
func (h *FavoriteHandler) validateFavoriteRequest(req *dto.FavoriteRequest) error {
	if req.BookID == 0 {
		return &ValidationError{Field: "book_id", Message: "Book ID is required"}
	}
	return nil
}
