package model

import "gorm.io/gorm"

// Favorite represents the database entity for user's favorite books
type Favorite struct {
	gorm.Model
	UserID uint `json:"user_id"`
	BookID uint `json:"book_id"`
}
