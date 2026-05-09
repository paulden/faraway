package repository

import (
	"github.com/paulden/faraway/models"
	"gorm.io/gorm"
)

type RoundRepository interface {
	FindAllByGameID(gameID uint) ([]models.Round, error)
	FindByID(id uint) (*models.Round, error)
	ExistsByNumberAndGameID(number uint, gameID uint) (bool, error)
	Create(round *models.Round) error
	Update(round *models.Round) error
	Delete(id uint) error
}

type roundRepository struct {
	db *gorm.DB
}

func NewRoundRepository(db *gorm.DB) RoundRepository {
	return &roundRepository{db: db}
}

func (r *roundRepository) FindAllByGameID(gameID uint) ([]models.Round, error) {
	var rounds []models.Round
	result := r.db.Where("game_id = ?", gameID).Order("number").Find(&rounds)
	return rounds, result.Error
}

func (r *roundRepository) FindByID(id uint) (*models.Round, error) {
	var round models.Round
	result := r.db.First(&round, id)
	return &round, result.Error
}

func (r *roundRepository) ExistsByNumberAndGameID(number uint, gameID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Round{}).Where("number = ? AND game_id = ?", number, gameID).Count(&count)
	return count > 0, result.Error
}

func (r *roundRepository) Create(round *models.Round) error {
	return r.db.Create(round).Error
}

func (r *roundRepository) Update(round *models.Round) error {
	return r.db.Save(round).Error
}

func (r *roundRepository) Delete(id uint) error {
	return r.db.Delete(&models.Round{}, id).Error
}
