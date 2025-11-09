package service

import (
	"bms-go/internal/infra/repository"
	"bms-go/internal/model"
	"bms-go/internal/model/dto"
	"fmt"
)

type FavoriteService struct {
	repo     *repository.FavoriteRepository
	bookRepo *repository.BookRepository
}

func NewFavoriteService(repo *repository.FavoriteRepository, bookRepo *repository.BookRepository) *FavoriteService {
	return &FavoriteService{repo: repo, bookRepo: bookRepo}
}

func (s *FavoriteService) GetFavorites(userID uint) ([]dto.FavoriteResponse, error) {
	favs, err := s.repo.FindAll(userID)
	if err != nil {
		return nil, err
	}

	var responses []dto.FavoriteResponse
	for _, f := range favs {
		book, err := s.bookRepo.FindByID(f.BookID)
		if err != nil {
			continue
		}

		responses = append(responses, dto.FavoriteResponse{
			ID:     f.ID,
			UserID: f.UserID,
			BookID: f.BookID,
			Book: &dto.BookResponse{
				ID:       book.ID,
				Title:    book.Title,
				Author:   book.Author,
				Category: book.Category,
			},
		})
	}

	return responses, nil
}

// GetFavoriteByID retrieves a single favorite by ID for a specific user
func (s *FavoriteService) GetFavoriteByID(userID, favoriteID uint) (*dto.FavoriteResponse, error) {
	fav, err := s.repo.FindByID(userID, favoriteID)
	if err != nil {
		return nil, err
	}

	book, err := s.bookRepo.FindByID(fav.BookID)
	if err != nil {
		return nil, err
	}

	return &dto.FavoriteResponse{
		ID:        fav.ID,
		UserID:    fav.UserID,
		BookID:    fav.BookID,
		CreatedAt: fav.CreatedAt,
		Book: &dto.BookResponse{
			ID:        book.ID,
			Title:     book.Title,
			Author:    book.Author,
			Category:  book.Category,
			CreatedAt: book.CreatedAt,
			UpdatedAt: book.UpdatedAt,
		},
	}, nil
}

func (s *FavoriteService) AddFavorite(userID uint, req dto.FavoriteRequest) (*dto.FavoriteResponse, error) {
	// Check if book exists first
	_, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, fmt.Errorf("book not found")
	}

	// Check if already favorited
	exists, err := s.repo.Exists(userID, req.BookID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("already in favorites")
	}

	fav := model.Favorite{
		UserID: userID,
		BookID: req.BookID,
	}

	if err := s.repo.Create(&fav); err != nil {
		return nil, err
	}

	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, err
	}

	return &dto.FavoriteResponse{
		ID:        fav.ID,
		UserID:    userID,
		BookID:    req.BookID,
		CreatedAt: fav.CreatedAt,
		Book: &dto.BookResponse{
			ID:        book.ID,
			Title:     book.Title,
			Author:    book.Author,
			Category:  book.Category,
			CreatedAt: book.CreatedAt,
			UpdatedAt: book.UpdatedAt,
		},
	}, nil
}

// RemoveFavorite deletes a favorite entry
func (s *FavoriteService) RemoveFavorite(userID, favoriteID uint) error {
	return s.repo.Delete(userID, favoriteID)
}
