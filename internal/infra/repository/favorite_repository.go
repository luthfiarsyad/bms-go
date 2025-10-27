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

func (r *FavoriteRepository) Create(fav *model.Favorite) error {
	return r.db.Create(fav).Error
}

func (r *FavoriteRepository) Delete(userID, favoriteID uint) error {
	return r.db.Where("id = ? AND user_id = ?", favoriteID, userID).Delete(&model.Favorite{}).Error
}
