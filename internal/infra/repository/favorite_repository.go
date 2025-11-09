package repository

import (
	"bms-go/internal/model"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) FindAll(userID uint) ([]model.Favorite, error) {
	var favs []model.Favorite
	if err := r.db.Preload("Book").Where("user_id = ?", userID).Find(&favs).Error; err != nil {
		return nil, err
	}
	return favs, nil
}

// FindByID retrieves a single favorite by ID for a specific user
func (r *FavoriteRepository) FindByID(userID, favoriteID uint) (*model.Favorite, error) {
	var fav model.Favorite
	if err := r.db.Where("id = ? AND user_id = ?", favoriteID, userID).First(&fav).Error; err != nil {
		return nil, err
	}
	return &fav, nil
}

// Exists checks if a favorite already exists for a user and book
func (r *FavoriteRepository) Exists(userID, bookID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ? AND book_id = ?", userID, bookID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FavoriteRepository) Create(fav *model.Favorite) error {
	return r.db.Create(fav).Error
}

func (r *FavoriteRepository) Delete(userID, favoriteID uint) error {
	return r.db.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(&model.Favorite{}).Error
}

// DeleteByBookID removes a favorite by user ID and book ID
func (r *FavoriteRepository) DeleteByBookID(userID, bookID uint) error {
	return r.db.Where("user_id = ? AND book_id = ?", userID, bookID).Delete(&model.Favorite{}).Error
}
