package service

import (
	"bms-go/internal/infra/repository"
	"bms-go/internal/model"
)

type BookService struct {
	repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) GetBooks(search, category string) ([]model.Book, error) {
	return s.repo.FindAll(search, category)
}

func (s *BookService) GetBookByID(id uint) (*model.Book, error) {
	return s.repo.FindByID(id)
}

func (s *BookService) CreateBook(book *model.Book) error {
	return s.repo.Create(book)
}

func (s *BookService) UpdateBook(book *model.Book) error {
	return s.repo.Update(book)
}

func (s *BookService) DeleteBook(id uint) error {
	return s.repo.Delete(id)
}
