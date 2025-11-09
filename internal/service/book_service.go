package service

import (
	"bms-go/internal/infra/repository"
	"bms-go/internal/model"
	"errors"
	"strings"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetBooks(search, category string) ([]model.Book, error) {
	search = strings.TrimSpace(search)
	category = strings.TrimSpace(category)
	
	return s.repo.FindAll(search, category)
}

// AdvancedSearch performs sophisticated search with multiple criteria
func (s *BookService) AdvancedSearch(params repository.AdvancedSearchParams) ([]model.Book, error) {
	// Validate and set defaults
	if params.SearchType == "" {
		params.SearchType = "contains"
	}
	if params.SortBy == "" {
		params.SortBy = "relevance"
	}
	if params.SortOrder == "" {
		params.SortOrder = "ASC"
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	
	// Validate search type
	validSearchTypes := map[string]bool{
		"exact":       true,
		"starts_with": true,
		"contains":    true,
		"fuzzy":       true,
	}
	if !validSearchTypes[params.SearchType] {
		return nil, errors.New("invalid search type. Must be: exact, starts_with, contains, or fuzzy")
	}
	
	// Validate sort field
	validSortFields := map[string]bool{
		"title":       true,
		"author":      true,
		"category":    true,
		"created_at":  true,
		"relevance":   true,
	}
	if !validSortFields[params.SortBy] {
		return nil, errors.New("invalid sort field. Must be: title, author, category, created_at, or relevance")
	}
	
	// Validate sort order
	if params.SortOrder != "ASC" && params.SortOrder != "DESC" {
		return nil, errors.New("invalid sort order. Must be: ASC or DESC")
	}
	
	return s.repo.AdvancedSearch(params)
}

// GetSearchSuggestions provides search suggestions
func (s *BookService) GetSearchSuggestions(query string, limit int) ([]string, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}
	
	return s.repo.GetSearchSuggestions(query, limit)
}

func (s *BookService) GetBookByID(id uint) (*model.Book, error) {
	if id == 0 {
		return nil, errors.New("invalid book ID")
	}
	
	return s.repo.FindByID(id)
}

func (s *BookService) CreateBook(book *model.Book) error {
	// Validate book data
	if err := s.validateBook(book); err != nil {
		return err
	}
	
	// Check for duplicate title
	existingBook, err := s.repo.FindByTitle(book.Title)
	if err == nil && existingBook != nil {
		return errors.New("book with this title already exists")
	}
	
	return s.repo.Create(book)
}

func (s *BookService) UpdateBook(book *model.Book) error {
	// Validate book data
	if err := s.validateBook(book); err != nil {
		return err
	}
	
	// Check if book exists
	exists, err := s.repo.Exists(book.ID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("book not found")
	}
	
	// Check for duplicate title (excluding current book)
	existingBook, err := s.repo.FindByTitle(book.Title)
	if err == nil && existingBook != nil && existingBook.ID != book.ID {
		return errors.New("book with this title already exists")
	}
	
	return s.repo.Update(book)
}

func (s *BookService) DeleteBook(id uint) error {
	if id == 0 {
		return errors.New("invalid book ID")
	}
	
	// Check if book exists
	exists, err := s.repo.Exists(id)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("book not found")
	}
	
	return s.repo.Delete(id)
}

// validateBook validates the book data
func (s *BookService) validateBook(book *model.Book) error {
	if strings.TrimSpace(book.Title) == "" {
		return errors.New("title is required")
	}
	if len(book.Title) > 255 {
		return errors.New("title must be less than 255 characters")
	}
	if strings.TrimSpace(book.Author) == "" {
		return errors.New("author is required")
	}
	if len(book.Author) > 255 {
		return errors.New("author must be less than 255 characters")
	}
	if strings.TrimSpace(book.Category) == "" {
		return errors.New("category is required")
	}
	if len(book.Category) > 255 {
		return errors.New("category must be less than 255 characters")
	}
	return nil
}

// GetBookCount returns the total number of books
func (s *BookService) GetBookCount() (int64, error) {
	return s.repo.GetCount()
}
