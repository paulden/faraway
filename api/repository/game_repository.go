package repository

import (
	"github.com/paulden/faraway/models"
	"gorm.io/gorm"
)

type GameRepository interface {
	FindAll() ([]models.Game, error)
	FindByID(id uint) (*models.Game, error)
	Create(game *models.Game) error
	Update(game *models.Game) error
	Delete(id uint) error
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) FindAll() ([]models.Game, error) {
	var games []models.Game
	result := r.db.Find(&games)
	return games, result.Error
}

func (r *gameRepository) FindByID(id uint) (*models.Game, error) {
	var game models.Game
	result := r.db.First(&game, id)
	return &game, result.Error
}

func (r *gameRepository) Create(game *models.Game) error {
	return r.db.Create(game).Error
}

func (r *gameRepository) Update(game *models.Game) error {
	return r.db.Save(game).Error
}

func (r *gameRepository) Delete(id uint) error {
	return r.db.Delete(&models.Game{}, id).Error
}
