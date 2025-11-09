package handler

import (
	"bms-go/internal/model"
	"bms-go/internal/model/dto"
	"bms-go/internal/infra/repository"
	"bms-go/internal/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	service *service.BookService
}

func NewBookHandler(s *service.BookService) *BookHandler {
	return &BookHandler{service: s}
}

func (h *BookHandler) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/books")
	group.GET("", h.GetBooks)
	group.GET("/search", h.AdvancedSearch)
	group.GET("/suggestions", h.GetSearchSuggestions)
	group.GET("/:id", h.GetBookByID)
	group.POST("", h.CreateBook)
	group.PUT("/:id", h.UpdateBook)
	group.DELETE("/:id", h.DeleteBook)
}

// GetBooks godoc
// @Summary Get all books with basic search
// @Description Retrieve a list of all books with optional basic search and category filtering. This endpoint provides simple search functionality that searches within book titles and authors.
// @Tags Books
// @Accept json
// @Produce json
// @Param search query string false "Search keyword to find books by title or author (basic contains search)"
// @Param category query string false "Filter books by category exactly"
// @Success 200 {object} dto.APIResponse{data=[]model.SwaggerBook} "Books retrieved successfully"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /books [get]
// @Example {
//   "request": "GET /books?search=harry&category=Fantasy",
//   "response": {
//     "success": true,
//     "message": "Books retrieved successfully",
//     "data": [
//       {
//         "id": 1,
//         "title": "Harry Potter and the Sorcerer's Stone",
//         "author": "J.K. Rowling",
//         "category": "Fantasy",
//         "created_at": "2023-01-01T00:00:00Z",
//         "updated_at": "2023-01-01T00:00:00Z"
//       }
//     ]
//   }
// }
func (h *BookHandler) GetBooks(c *gin.Context) {
	search := strings.TrimSpace(c.Query("search"))
	category := strings.TrimSpace(c.Query("category"))

	books, err := h.service.GetBooks(search, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to retrieve books",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Data:    books,
	})
}

// GetBookByID godoc
// @Summary Get book by ID
// @Description Retrieve detailed information about a specific book using its unique identifier. Returns complete book details including timestamps.
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID (must be a positive integer)"
// @Success 200 {object} dto.APIResponse{data=model.SwaggerBook} "Book retrieved successfully"
// @Failure 400 {object} dto.APIResponse "Invalid book ID format"
// @Failure 404 {object} dto.APIResponse "Book not found"
// @Router /books/{id} [get]
// @Example {
//   "request": "GET /books/1",
//   "response": {
//     "success": true,
//     "message": "Book retrieved successfully",
//     "data": {
//       "id": 1,
//       "title": "The Great Gatsby",
//       "author": "F. Scott Fitzgerald",
//       "category": "Classic",
//       "created_at": "2023-01-01T00:00:00Z",
//       "updated_at": "2023-01-01T00:00:00Z"
//     }
//   }
// }
func (h *BookHandler) GetBookByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   "Book ID must be a positive integer",
		})
		return
	}

	book, err := h.service.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

// CreateBook godoc
// @Summary Create new book
// @Description Add a new book to the system. The system validates for duplicate titles and ensures all required fields are present. Title, author, and category must be between 1-255 characters.
// @Tags Books
// @Accept json
// @Produce json
// @Param book body dto.BookRequest true "Book information" required(true)
// @Success 201 {object} dto.APIResponse{data=model.Book} "Book created successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request body or validation failed"
// @Failure 409 {object} dto.APIResponse "Book with this title already exists"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /books [post]
// @Example {
//   "request": {
//     "title": "1984",
//     "author": "George Orwell",
//     "category": "Dystopian"
//   },
//   "response": {
//     "success": true,
//     "message": "Book created successfully",
//     "data": {
//       "id": 1,
//       "title": "1984",
//       "author": "George Orwell",
//       "category": "Dystopian",
//       "created_at": "2023-01-01T00:00:00Z",
//       "updated_at": "2023-01-01T00:00:00Z"
//     }
//   }
// }
func (h *BookHandler) CreateBook(c *gin.Context) {
	var req dto.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate input
	if err := h.validateBookRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	book := &model.Book{
		Title:    strings.TrimSpace(req.Title),
		Author:   strings.TrimSpace(req.Author),
		Category: strings.TrimSpace(req.Category),
	}

	if err := h.service.CreateBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to create book",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Book created successfully",
		Data:    book,
	})
}

// UpdateBook godoc
// @Summary Update book
// @Description Update existing book information by ID. Validates that the book exists and checks for duplicate titles (excluding the current book). All fields are required and must be between 1-255 characters.
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID (must be a positive integer)"
// @Param book body model.SwaggerBookRequest true "Updated book information" required(true)
// @Success 200 {object} dto.APIResponse{data=model.SwaggerBook} "Book updated successfully"
// @Failure 400 {object} dto.APIResponse "Invalid request body, validation failed, or invalid book ID"
// @Failure 404 {object} dto.APIResponse "Book not found"
// @Failure 409 {object} dto.APIResponse "Book with this title already exists"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /books/{id} [put]
// @Example {
//   "request": "PUT /books/1",
//   "body": {
//     "title": "1984: Special Edition",
//     "author": "George Orwell",
//     "category": "Classic Dystopian"
//   },
//   "response": {
//     "success": true,
//     "message": "Book updated successfully",
//     "data": {
//       "id": 1,
//       "title": "1984: Special Edition",
//       "author": "George Orwell",
//       "category": "Classic Dystopian",
//       "created_at": "2023-01-01T00:00:00Z",
//       "updated_at": "2023-01-01T01:00:00Z"
//     }
//   }
// }
func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   "Book ID must be a positive integer",
		})
		return
	}

	var req dto.BookRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate input
	if err := h.validateBookRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Check if book exists first
	_, err = h.service.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   err.Error(),
		})
		return
	}

	book := &model.Book{
		Title:    strings.TrimSpace(req.Title),
		Author:   strings.TrimSpace(req.Author),
		Category: strings.TrimSpace(req.Category),
	}
	book.ID = uint(id)

	if err := h.service.UpdateBook(book); err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to update book",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Book updated successfully",
		Data:    book,
	})
}

// DeleteBook godoc
// @Summary Delete book
// @Description Soft delete a book by its ID. The book is marked as deleted but remains in the database. This operation is irreversible and will also remove the book from all favorites lists.
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID (must be a positive integer)"
// @Success 200 {object} dto.APIResponse "Book deleted successfully"
// @Failure 400 {object} dto.APIResponse "Invalid book ID format"
// @Failure 404 {object} dto.APIResponse "Book not found"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /books/{id} [delete]
// @Example {
//   "request": "DELETE /books/1",
//   "response": {
//     "success": true,
//     "message": "Book deleted successfully"
//   }
// }
func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   "Book ID must be a positive integer",
		})
		return
	}

	// Check if book exists first
	_, err = h.service.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   err.Error(),
		})
		return
	}

	if err := h.service.DeleteBook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to delete book",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Book deleted successfully",
	})
}

// validateBookRequest validates the book request data
func (h *BookHandler) validateBookRequest(req *dto.BookRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return &ValidationError{Field: "title", Message: "Title is required"}
	}
	if len(req.Title) > 255 {
		return &ValidationError{Field: "title", Message: "Title must be less than 255 characters"}
	}
	if strings.TrimSpace(req.Author) == "" {
		return &ValidationError{Field: "author", Message: "Author is required"}
	}
	if len(req.Author) > 255 {
		return &ValidationError{Field: "author", Message: "Author must be less than 255 characters"}
	}
	if strings.TrimSpace(req.Category) == "" {
		return &ValidationError{Field: "category", Message: "Category is required"}
	}
	if len(req.Category) > 255 {
		return &ValidationError{Field: "category", Message: "Category must be less than 255 characters"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// AdvancedSearch godoc
// @Summary Advanced book search with relevance scoring
// @Description Perform sophisticated search with multiple criteria, relevance scoring, and advanced matching algorithms. Supports exact matching, prefix matching, contains search, and fuzzy search with typo tolerance. Results are ranked by relevance with exact matches appearing first.
// @Tags Books
// @Accept json
// @Produce json
// @Param query query string false "Search query string - searches across title, author, and category"
// @Param category query string false "Filter by exact category match"
// @Param author query string false "Filter by author (partial match)"
// @Param search_type query string false "Search strategy" Enums(exact, starts_with, contains, fuzzy) default("contains")
// @Param sort_by query string false "Sort field for results" Enums(title, author, category, created_at, relevance) default("relevance")
// @Param sort_order query string false "Sort order" Enums(ASC, DESC) default("ASC")
// @Param limit query int false "Maximum number of results to return (1-100)" default(20)
// @Param offset query int false "Number of results to skip for pagination" default(0)
// @Success 200 {object} dto.APIResponse{data=[]model.SwaggerBook} "Search completed successfully"
// @Failure 400 {object} dto.APIResponse "Invalid search parameters"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /books/search [get]
// @Example {
//   "request": "GET /books/search?query=harry%20potter&search_type=contains&sort_by=relevance&limit=10",
//   "response": {
//     "success": true,
//     "message": "Search completed successfully",
//     "data": [
//       {
//         "id": 1,
//         "title": "Harry Potter and the Sorcerer's Stone",
//         "author": "J.K. Rowling",
//         "category": "Fantasy",
//         "created_at": "2023-01-01T00:00:00Z",
//         "updated_at": "2023-01-01T00:00:00Z"
//       }
//     ]
//   }
// }
func (h *BookHandler) AdvancedSearch(c *gin.Context) {
	params := repository.AdvancedSearchParams{
		Query:      strings.TrimSpace(c.Query("query")),
		Category:   strings.TrimSpace(c.Query("category")),
		Author:     strings.TrimSpace(c.Query("author")),
		SearchType: strings.TrimSpace(c.Query("search_type")),
		SortBy:     strings.TrimSpace(c.Query("sort_by")),
		SortOrder:  strings.TrimSpace(c.Query("sort_order")),
	}

	// Parse limit and offset
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			params.Offset = offset
		}
	}

	books, err := h.service.AdvancedSearch(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Search failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Search completed successfully",
		Data:    books,
	})
}

// GetSearchSuggestions godoc
// @Summary Get search suggestions for autocomplete
// @Description Retrieve search suggestions based on existing book titles and authors. Perfect for implementing autocomplete functionality in user interfaces. Returns unique titles and authors that contain the search query.
// @Tags Books
// @Accept json
// @Produce json
// @Param query query string true "Search query for suggestions (minimum 1 character)"
// @Param limit query int false "Maximum number of suggestions to return (1-20)" default(10)
// @Success 200 {object} dto.APIResponse{data=[]string} "Suggestions retrieved successfully"
// @Failure 400 {object} dto.APIResponse "Query parameter is required or invalid limit"
// @Failure 500 {object} dto.APIResponse "Internal server error"
// @Router /books/suggestions [get]
// @Example {
//   "request": "GET /books/suggestions?query=harry&limit=5",
//   "response": {
//     "success": true,
//     "message": "Suggestions retrieved successfully",
//     "data": [
//       "Harry Potter and the Sorcerer's Stone",
//       "Harry Potter and the Chamber of Secrets",
//       "J.K. Rowling"
//     ]
//   }
// }
func (h *BookHandler) GetSearchSuggestions(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Query parameter is required",
			Error:   "Search query cannot be empty",
		})
		return
	}

	limit := 10
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	suggestions, err := h.service.GetSearchSuggestions(query, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: "Failed to get suggestions",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Suggestions retrieved successfully",
		Data:    suggestions,
	})
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
