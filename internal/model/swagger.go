package model

import "time"

// SwaggerBook represents a book for Swagger documentation (without GORM dependencies)
type SwaggerBook struct {
	ID        uint      `json:"id" example:"1"`
	Title     string    `json:"title" example:"Harry Potter and the Sorcerer's Stone"`
	Author    string    `json:"author" example:"J.K. Rowling"`
	Category  string    `json:"category" example:"Fantasy"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SwaggerFavorite represents a favorite for Swagger documentation
type SwaggerFavorite struct {
	ID        uint              `json:"id" example:"1"`
	UserID    uint              `json:"user_id" example:"1"`
	BookID    uint              `json:"book_id" example:"1"`
	CreatedAt time.Time         `json:"created_at" example:"2023-01-01T00:00:00Z"`
	Book      *SwaggerBook      `json:"book,omitempty"`
}

// SwaggerBookRequest represents a book request for Swagger documentation
type SwaggerBookRequest struct {
	Title    string `json:"title" example:"1984" binding:"required"`
	Author   string `json:"author" example:"George Orwell" binding:"required"`
	Category string `json:"category" example:"Dystopian" binding:"required"`
}

// SwaggerFavoriteRequest represents a favorite request for Swagger documentation
type SwaggerFavoriteRequest struct {
	BookID uint `json:"book_id" example:"1" binding:"required"`
}