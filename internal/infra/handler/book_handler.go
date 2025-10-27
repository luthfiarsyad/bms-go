package handler

import (
	"bms-go/internal/model"
	"bms-go/internal/service"
	"net/http"
	"strconv"

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
	group.GET("/:id", h.GetBookByID)
	group.POST("", h.CreateBook)
	group.PUT("/:id", h.UpdateBook)
	group.DELETE("/:id", h.DeleteBook)
}

// GetBooks godoc
// @Summary Get all books
// @Description Get list of all books, optionally filtered by search or category
// @Tags Books
// @Accept json
// @Produce json
// @Param search query string false "Search keyword"
// @Param category query string false "Category filter"
// @Success 200 {array} model.Book
// @Failure 500 {object} map[string]string
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	search := c.Query("search")
	category := c.Query("category")

	books, err := h.service.GetBooks(search, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetBookByID godoc
// @Summary Get book by ID
// @Description Retrieve a single book by its ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} model.Book
// @Failure 404 {object} map[string]string
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	book, err := h.service.GetBookByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// CreateBook godoc
// @Summary Create new book
// @Description Add a new book to the system
// @Tags Books
// @Accept json
// @Produce json
// @Param book body model.Book true "Book object"
// @Success 201 {object} model.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// UpdateBook godoc
// @Summary Update book
// @Description Update book information by ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body model.Book true "Updated book data"
// @Success 200 {object} model.Book
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	book.ID = uint(id)
	if err := h.service.UpdateBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DeleteBook godoc
// @Summary Delete book
// @Description Delete a book by its ID
// @Tags Books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.service.DeleteBook(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
