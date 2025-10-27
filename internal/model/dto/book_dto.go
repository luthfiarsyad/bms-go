package dto

type BookRequest struct {
	Title    string `json:"title" binding:"required"`
	Author   string `json:"author" binding:"required"`
	Category string `json:"category" binding:"required"`
}

type BookResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Category string `json:"category"`
}
