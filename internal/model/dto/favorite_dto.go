package dto

import "time"

// FavoriteRequest represents the request payload for adding a favorite
type FavoriteRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

// FavoriteResponse represents the response payload for a favorite
type FavoriteResponse struct {
	ID        uint         `json:"id"`
	UserID    uint         `json:"user_id"`
	BookID    uint         `json:"book_id"`
	CreatedAt time.Time    `json:"created_at"`
	Book      *BookResponse `json:"book,omitempty"`
}

// FavoriteListResponse represents the response payload for a list of favorites
type FavoriteListResponse struct {
	Favorites []FavoriteResponse `json:"favorites"`
	Count     int                `json:"count"`
}
