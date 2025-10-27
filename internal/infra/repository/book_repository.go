package repository

import (
	"bms-go/internal/model"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindAll(search, category string) ([]model.Book, error) {
	var books []model.Book
	query := r.db

	if search != "" {
		query = query.Where("title LIKE ? OR author LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (r *BookRepository) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	if err := r.db.First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Save(book).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&model.Book{}, id).Error
}