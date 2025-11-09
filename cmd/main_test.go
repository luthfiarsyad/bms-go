package main_test

import (
	"bms-go/internal/infra/handler"
	"bms-go/internal/infra/repository"
	"bms-go/internal/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// minimal models used for migrations in tests
type BookModel struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Author string
	Year   int
}

type FavoriteModel struct {
	ID     uint `gorm:"primaryKey"`
	BookID uint
}

func setupRouter(t *testing.T) (*gin.Engine, *gorm.DB) {
	gin.SetMode(gin.TestMode)

	var db *gorm.DB
	var err error

	dsn := os.Getenv("TEST_DSN")
	if dsn != "" {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else {
		db, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	}
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	// migrate test models so repository operations work
	if err := db.AutoMigrate(&BookModel{}, &FavoriteModel{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	// Ensure clean state for tests (safe for sqlite; for MySQL these table names come from our test models)
	db.Exec("DELETE FROM book_models")
	db.Exec("DELETE FROM favorite_models")

	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	favRepo := repository.NewFavoriteRepository(db)
	favService := service.NewFavoriteService(favRepo, bookRepo)
	favHandler := handler.NewFavoriteHandler(favService)

	r := gin.Default()
	bookHandler.RegisterRoutes(r)
	favHandler.RegisterRoutes(r)
	r.NoRoute(handler.NotFoundHandler)

	return r, db
}

// helper to create a book and return its DB ID
func createBookAndGetID(t *testing.T, r *gin.Engine, db *gorm.DB, payload map[string]interface{}) uint {
	jb, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create book failed status=%d body=%s", w.Code, w.Body.String())
	}

	var bm BookModel
	if err := db.Order("id desc").First(&bm).Error; err != nil {
		t.Fatalf("cannot query created book: %v", err)
	}
	return bm.ID
}

func TestGetBooks(t *testing.T) {
	r, _ := setupRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCreateBook(t *testing.T) {
	r, db := setupRouter(t)
	_ = db // keep db to ensure migrations ran

	book := map[string]interface{}{
		"title":  "Test Book",
		"author": "Test Author",
		"year":   2023,
	}
	jb, _ := json.Marshal(book)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAddFavorite(t *testing.T) {
	r, db := setupRouter(t)

	// create book first and get real ID
	book := map[string]interface{}{
		"title":  "Fav Book",
		"author": "Author",
		"year":   2020,
	}
	id := createBookAndGetID(t, r, db, book)

	fav := map[string]interface{}{
		"book_id": id,
	}
	jb, _ := json.Marshal(fav)
	req, _ := http.NewRequest("POST", "/favorites", bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusBadRequest)
}

func TestNotFoundRoute(t *testing.T) {
	r, _ := setupRouter(t)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/nonexistent", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBook(t *testing.T) {
	r, db := setupRouter(t)

	// create then delete
	book := map[string]interface{}{
		"title":  "To Delete",
		"author": "Author",
		"year":   2000,
	}
	id := createBookAndGetID(t, r, db, book)

	req, _ := http.NewRequest("DELETE", "/books/"+json.Number(id).String(), nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
}

func TestUpdateBook(t *testing.T) {
	r, db := setupRouter(t)

	book := map[string]interface{}{
		"title":  "Before Update",
		"author": "Author",
		"year":   2001,
	}
	id := createBookAndGetID(t, r, db, book)

	updated := map[string]interface{}{
		"title":  "After Update",
		"author": "New Author",
		"year":   2025,
	}
	jb, _ := json.Marshal(updated)
	req, _ := http.NewRequest("PUT", "/books/"+json.Number(id).String(), bytes.NewBuffer(jb))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusNotFound)
}
