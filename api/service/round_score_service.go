package service

import (
	"errors"

	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/repository"
)

var (
	ErrRoundScoreNotFound = errors.New("score not found")
	ErrPlayerNotInGame    = errors.New("player does not belong to this game")
	ErrScoreAlreadyExists = errors.New("this player already has a score for this round")
)

type RoundScoreService interface {
	GetAllByRound(roundID uint) ([]models.RoundScore, error)
	AddScore(roundID uint, score *models.RoundScore) (*models.RoundScore, error)
	UpdateScore(id uint, newScore int) (*models.RoundScore, error)
	DeleteScore(id uint) error
}

type roundScoreService struct {
	scoreRepo  repository.RoundScoreRepository
	roundRepo  repository.RoundRepository
	playerRepo repository.PlayerRepository
}

func NewRoundScoreService(
	scoreRepo repository.RoundScoreRepository,
	roundRepo repository.RoundRepository,
	playerRepo repository.PlayerRepository,
) RoundScoreService {
	return &roundScoreService{
		scoreRepo:  scoreRepo,
		roundRepo:  roundRepo,
		playerRepo: playerRepo,
	}
}

func (s *roundScoreService) GetAllByRound(roundID uint) ([]models.RoundScore, error) {
	if _, err := s.roundRepo.FindByID(roundID); err != nil {
		return nil, ErrRoundNotFound
	}
	return s.scoreRepo.FindAllByRoundID(roundID)
}

func (s *roundScoreService) AddScore(roundID uint, score *models.RoundScore) (*models.RoundScore, error) {
	round, err := s.roundRepo.FindByID(roundID)
	if err != nil {
		return nil, ErrRoundNotFound
	}

	player, err := s.playerRepo.FindByID(score.PlayerID)
	if err != nil {
		return nil, errors.New("player not found")
	}
	if player.GameID != round.GameID {
		return nil, ErrPlayerNotInGame
	}

	exists, err := s.scoreRepo.ExistsByPlayerAndRoundID(score.PlayerID, roundID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrScoreAlreadyExists
	}

	score.RoundID = roundID
	if err := s.scoreRepo.Create(score); err != nil {
		return nil, err
	}
	return score, nil
}

func (s *roundScoreService) UpdateScore(id uint, newScore int) (*models.RoundScore, error) {
	score, err := s.scoreRepo.FindByID(id)
	if err != nil {
		return nil, ErrRoundScoreNotFound
	}

	score.Score = newScore
	if err := s.scoreRepo.Update(score); err != nil {
		return nil, err
	}
	return score, nil
}

func (s *roundScoreService) DeleteScore(id uint) error {
	if _, err := s.scoreRepo.FindByID(id); err != nil {
		return ErrRoundScoreNotFound
	}
	return s.scoreRepo.Delete(id)
}
