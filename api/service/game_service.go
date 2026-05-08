package service

import (
	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/repository"
)

type GameService interface {
	GetAll() ([]models.Game, error)
	GetByID(id uint) (*models.Game, error)
	Create(game *models.Game) error
	Update(id uint, input *models.Game) (*models.Game, error)
	Delete(id uint) error
}

type gameService struct {
	repo repository.GameRepository
}

func NewGameService(repo repository.GameRepository) GameService {
	return &gameService{repo: repo}
}

func (s *gameService) GetAll() ([]models.Game, error) {
	return s.repo.FindAll()
}

func (s *gameService) GetByID(id uint) (*models.Game, error) {
	return s.repo.FindByID(id)
}

func (s *gameService) Create(game *models.Game) error {
	return s.repo.Create(game)
}

func (s *gameService) Update(id uint, input *models.Game) (*models.Game, error) {
	game, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	game.Title = input.Title
	game.IsFinished = input.IsFinished

	if err := s.repo.Update(game); err != nil {
		return nil, err
	}
	return game, nil
}

func (s *gameService) Delete(id uint) error {
	if _, err := s.repo.FindByID(id); err != nil {
		return err
	}
	return s.repo.Delete(id)
}
