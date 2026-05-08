package service

import (
	"errors"

	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/repository"
)

const MaxPlayersPerGame = 8

var (
	ErrGameNotFound       = errors.New("game not found")
	ErrDuplicatePlayerName = errors.New("a player with this name already exists in this game")
	ErrMaxPlayersReached  = errors.New("this game already has the maximum number of players (8)")
)

type PlayerService interface {
	GetAllByGame(gameID uint) ([]models.Player, error)
	AddPlayer(gameID uint, player *models.Player) (*models.Player, error)
	RemovePlayer(id uint) error
}

type playerService struct {
	playerRepo repository.PlayerRepository
	gameRepo   repository.GameRepository
}

func NewPlayerService(playerRepo repository.PlayerRepository, gameRepo repository.GameRepository) PlayerService {
	return &playerService{playerRepo: playerRepo, gameRepo: gameRepo}
}

func (s *playerService) GetAllByGame(gameID uint) ([]models.Player, error) {
	if _, err := s.gameRepo.FindByID(gameID); err != nil {
		return nil, ErrGameNotFound
	}
	return s.playerRepo.FindAllByGameID(gameID)
}

func (s *playerService) AddPlayer(gameID uint, player *models.Player) (*models.Player, error) {
	if _, err := s.gameRepo.FindByID(gameID); err != nil {
		return nil, ErrGameNotFound
	}

	count, err := s.playerRepo.CountByGameID(gameID)
	if err != nil {
		return nil, err
	}
	if count >= MaxPlayersPerGame {
		return nil, ErrMaxPlayersReached
	}

	exists, err := s.playerRepo.ExistsByNameAndGameID(player.Name, gameID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicatePlayerName
	}

	player.GameID = gameID
	if err := s.playerRepo.Create(player); err != nil {
		return nil, err
	}
	return player, nil
}

func (s *playerService) RemovePlayer(id uint) error {
	if _, err := s.playerRepo.FindByID(id); err != nil {
		return errors.New("player not found")
	}
	return s.playerRepo.Delete(id)
}
