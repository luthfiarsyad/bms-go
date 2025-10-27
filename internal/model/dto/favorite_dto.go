package dto

type FavoriteRequest struct {
	BookID uint `json:"book_id" binding:"required"`
}

type FavoriteResponse struct {
	ID     uint          `json:"id"`
	UserID uint          `json:"user_id"`
	BookID uint          `json:"book_id"`
	Book   *BookResponse `json:"book,omitempty"` 
}
