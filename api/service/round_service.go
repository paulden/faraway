package service

import (
	"errors"

	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/repository"
)

var (
	ErrRoundNotFound    = errors.New("round not found")
	ErrRoundNumberExists = errors.New("a round with this number already exists in this game")
)

type RoundService interface {
	GetAllByGame(gameID uint) ([]models.Round, error)
	GetByID(id uint) (*models.Round, error)
	Create(gameID uint, round *models.Round) (*models.Round, error)
	Update(id uint, input *models.Round) (*models.Round, error)
	Delete(id uint) error
}

type roundService struct {
	roundRepo repository.RoundRepository
	gameRepo  repository.GameRepository
}

func NewRoundService(roundRepo repository.RoundRepository, gameRepo repository.GameRepository) RoundService {
	return &roundService{roundRepo: roundRepo, gameRepo: gameRepo}
}

func (s *roundService) GetAllByGame(gameID uint) ([]models.Round, error) {
	if _, err := s.gameRepo.FindByID(gameID); err != nil {
		return nil, ErrGameNotFound
	}
	return s.roundRepo.FindAllByGameID(gameID)
}

func (s *roundService) GetByID(id uint) (*models.Round, error) {
	round, err := s.roundRepo.FindByID(id)
	if err != nil {
		return nil, ErrRoundNotFound
	}
	return round, nil
}

func (s *roundService) Create(gameID uint, round *models.Round) (*models.Round, error) {
	if _, err := s.gameRepo.FindByID(gameID); err != nil {
		return nil, ErrGameNotFound
	}

	exists, err := s.roundRepo.ExistsByNumberAndGameID(round.Number, gameID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrRoundNumberExists
	}

	round.GameID = gameID
	if err := s.roundRepo.Create(round); err != nil {
		return nil, err
	}
	return round, nil
}

func (s *roundService) Update(id uint, input *models.Round) (*models.Round, error) {
	round, err := s.roundRepo.FindByID(id)
	if err != nil {
		return nil, ErrRoundNotFound
	}

	exists, err := s.roundRepo.ExistsByNumberAndGameID(input.Number, round.GameID)
	if err != nil {
		return nil, err
	}
	if exists && input.Number != round.Number {
		return nil, ErrRoundNumberExists
	}

	round.Number = input.Number
	if err := s.roundRepo.Update(round); err != nil {
		return nil, err
	}
	return round, nil
}

func (s *roundService) Delete(id uint) error {
	if _, err := s.roundRepo.FindByID(id); err != nil {
		return ErrRoundNotFound
	}
	return s.roundRepo.Delete(id)
}
