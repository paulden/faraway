package service

import (
	"errors"
	"testing"

	"github.com/paulden/faraway/models"
	"github.com/stretchr/testify/assert"
)

type mockPlayerRepository struct {
	findAllByGameID      func(gameID uint) ([]models.Player, error)
	findByID             func(id uint) (*models.Player, error)
	existsByNameAndGameID func(name string, gameID uint) (bool, error)
	countByGameID        func(gameID uint) (int64, error)
	create               func(player *models.Player) error
	delete               func(id uint) error
}

func (m *mockPlayerRepository) FindAllByGameID(gameID uint) ([]models.Player, error) {
	return m.findAllByGameID(gameID)
}
func (m *mockPlayerRepository) FindByID(id uint) (*models.Player, error) {
	return m.findByID(id)
}
func (m *mockPlayerRepository) ExistsByNameAndGameID(name string, gameID uint) (bool, error) {
	return m.existsByNameAndGameID(name, gameID)
}
func (m *mockPlayerRepository) CountByGameID(gameID uint) (int64, error) {
	return m.countByGameID(gameID)
}
func (m *mockPlayerRepository) Create(player *models.Player) error { return m.create(player) }
func (m *mockPlayerRepository) Delete(id uint) error               { return m.delete(id) }

type mockGameRepositoryForPlayer struct {
	findByID func(id uint) (*models.Game, error)
}

func (m *mockGameRepositoryForPlayer) FindAll() ([]models.Game, error)       { return nil, nil }
func (m *mockGameRepositoryForPlayer) FindByID(id uint) (*models.Game, error) { return m.findByID(id) }
func (m *mockGameRepositoryForPlayer) Create(game *models.Game) error         { return nil }
func (m *mockGameRepositoryForPlayer) Update(game *models.Game) error         { return nil }
func (m *mockGameRepositoryForPlayer) Delete(id uint) error                   { return nil }

var gameExists = &mockGameRepositoryForPlayer{
	findByID: func(id uint) (*models.Game, error) { return &models.Game{}, nil },
}
var gameNotFound = &mockGameRepositoryForPlayer{
	findByID: func(id uint) (*models.Game, error) { return nil, errors.New("record not found") },
}

// --- GetAllByGame ---

func TestGetAllByGame_Success(t *testing.T) {
	expected := []models.Player{{Name: "Alice"}}
	svc := NewPlayerService(&mockPlayerRepository{
		findAllByGameID: func(gameID uint) ([]models.Player, error) { return expected, nil },
	}, gameExists)

	players, err := svc.GetAllByGame(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, players)
}

func TestGetAllByGame_GameNotFound(t *testing.T) {
	svc := NewPlayerService(&mockPlayerRepository{}, gameNotFound)

	players, err := svc.GetAllByGame(99)

	assert.ErrorIs(t, err, ErrGameNotFound)
	assert.Nil(t, players)
}

// --- AddPlayer ---

func TestAddPlayer_Success(t *testing.T) {
	player := &models.Player{Name: "Alice"}
	svc := NewPlayerService(&mockPlayerRepository{
		countByGameID:         func(gameID uint) (int64, error) { return 3, nil },
		existsByNameAndGameID: func(name string, gameID uint) (bool, error) { return false, nil },
		create:                func(p *models.Player) error { return nil },
	}, gameExists)

	created, err := svc.AddPlayer(1, player)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), created.GameID)
}

func TestAddPlayer_GameNotFound(t *testing.T) {
	svc := NewPlayerService(&mockPlayerRepository{}, gameNotFound)

	created, err := svc.AddPlayer(99, &models.Player{Name: "Alice"})

	assert.ErrorIs(t, err, ErrGameNotFound)
	assert.Nil(t, created)
}

func TestAddPlayer_MaxPlayersReached(t *testing.T) {
	svc := NewPlayerService(&mockPlayerRepository{
		countByGameID: func(gameID uint) (int64, error) { return MaxPlayersPerGame, nil },
	}, gameExists)

	created, err := svc.AddPlayer(1, &models.Player{Name: "Alice"})

	assert.ErrorIs(t, err, ErrMaxPlayersReached)
	assert.Nil(t, created)
}

func TestAddPlayer_DuplicateName(t *testing.T) {
	svc := NewPlayerService(&mockPlayerRepository{
		countByGameID:         func(gameID uint) (int64, error) { return 2, nil },
		existsByNameAndGameID: func(name string, gameID uint) (bool, error) { return true, nil },
	}, gameExists)

	created, err := svc.AddPlayer(1, &models.Player{Name: "Alice"})

	assert.ErrorIs(t, err, ErrDuplicatePlayerName)
	assert.Nil(t, created)
}

// --- RemovePlayer ---

func TestRemovePlayer_Success(t *testing.T) {
	svc := NewPlayerService(&mockPlayerRepository{
		findByID: func(id uint) (*models.Player, error) { return &models.Player{}, nil },
		delete:   func(id uint) error { return nil },
	}, gameExists)

	err := svc.RemovePlayer(1)

	assert.NoError(t, err)
}

func TestRemovePlayer_NotFound(t *testing.T) {
	svc := NewPlayerService(&mockPlayerRepository{
		findByID: func(id uint) (*models.Player, error) { return nil, errors.New("record not found") },
	}, gameExists)

	err := svc.RemovePlayer(99)

	assert.Error(t, err)
}
