package repository

import (
	"github.com/paulden/faraway/models"
	"gorm.io/gorm"
)

type RoundScoreRepository interface {
	FindAllByRoundID(roundID uint) ([]models.RoundScore, error)
	FindByID(id uint) (*models.RoundScore, error)
	ExistsByPlayerAndRoundID(playerID uint, roundID uint) (bool, error)
	Create(score *models.RoundScore) error
	Update(score *models.RoundScore) error
	Delete(id uint) error
}

type roundScoreRepository struct {
	db *gorm.DB
}

func NewRoundScoreRepository(db *gorm.DB) RoundScoreRepository {
	return &roundScoreRepository{db: db}
}

func (r *roundScoreRepository) FindAllByRoundID(roundID uint) ([]models.RoundScore, error) {
	var scores []models.RoundScore
	result := r.db.Where("round_id = ?", roundID).Find(&scores)
	return scores, result.Error
}

func (r *roundScoreRepository) FindByID(id uint) (*models.RoundScore, error) {
	var score models.RoundScore
	result := r.db.First(&score, id)
	return &score, result.Error
}

func (r *roundScoreRepository) ExistsByPlayerAndRoundID(playerID uint, roundID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.RoundScore{}).Where("player_id = ? AND round_id = ?", playerID, roundID).Count(&count)
	return count > 0, result.Error
}

func (r *roundScoreRepository) Create(score *models.RoundScore) error {
	return r.db.Create(score).Error
}

func (r *roundScoreRepository) Update(score *models.RoundScore) error {
	return r.db.Save(score).Error
}

func (r *roundScoreRepository) Delete(id uint) error {
	return r.db.Delete(&models.RoundScore{}, id).Error
}
