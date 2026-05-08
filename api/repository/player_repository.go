package repository

import (
	"github.com/paulden/faraway/models"
	"gorm.io/gorm"
)

type PlayerRepository interface {
	FindAllByGameID(gameID uint) ([]models.Player, error)
	FindByID(id uint) (*models.Player, error)
	ExistsByNameAndGameID(name string, gameID uint) (bool, error)
	CountByGameID(gameID uint) (int64, error)
	Create(player *models.Player) error
	Delete(id uint) error
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db: db}
}

func (r *playerRepository) FindAllByGameID(gameID uint) ([]models.Player, error) {
	var players []models.Player
	result := r.db.Where("game_id = ?", gameID).Find(&players)
	return players, result.Error
}

func (r *playerRepository) FindByID(id uint) (*models.Player, error) {
	var player models.Player
	result := r.db.First(&player, id)
	return &player, result.Error
}

func (r *playerRepository) ExistsByNameAndGameID(name string, gameID uint) (bool, error) {
	var count int64
	result := r.db.Model(&models.Player{}).Where("name = ? AND game_id = ?", name, gameID).Count(&count)
	return count > 0, result.Error
}

func (r *playerRepository) CountByGameID(gameID uint) (int64, error) {
	var count int64
	result := r.db.Model(&models.Player{}).Where("game_id = ?", gameID).Count(&count)
	return count, result.Error
}

func (r *playerRepository) Create(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) Delete(id uint) error {
	return r.db.Delete(&models.Player{}, id).Error
}
