package service

import (
	"errors"
	"testing"

	"github.com/paulden/faraway/models"
	"github.com/stretchr/testify/assert"
)

type mockRoundScoreRepository struct {
	findAllByRoundID          func(roundID uint) ([]models.RoundScore, error)
	findByID                  func(id uint) (*models.RoundScore, error)
	existsByPlayerAndRoundID  func(playerID uint, roundID uint) (bool, error)
	create                    func(score *models.RoundScore) error
	update                    func(score *models.RoundScore) error
	delete                    func(id uint) error
}

func (m *mockRoundScoreRepository) FindAllByRoundID(roundID uint) ([]models.RoundScore, error) {
	return m.findAllByRoundID(roundID)
}
func (m *mockRoundScoreRepository) FindByID(id uint) (*models.RoundScore, error) {
	return m.findByID(id)
}
func (m *mockRoundScoreRepository) ExistsByPlayerAndRoundID(playerID uint, roundID uint) (bool, error) {
	return m.existsByPlayerAndRoundID(playerID, roundID)
}
func (m *mockRoundScoreRepository) Create(score *models.RoundScore) error { return m.create(score) }
func (m *mockRoundScoreRepository) Update(score *models.RoundScore) error { return m.update(score) }
func (m *mockRoundScoreRepository) Delete(id uint) error                  { return m.delete(id) }

type mockRoundRepoForScore struct {
	findByID func(id uint) (*models.Round, error)
}

func (m *mockRoundRepoForScore) FindAllByGameID(gameID uint) ([]models.Round, error) {
	return nil, nil
}
func (m *mockRoundRepoForScore) FindByID(id uint) (*models.Round, error) { return m.findByID(id) }
func (m *mockRoundRepoForScore) ExistsByNumberAndGameID(number uint, gameID uint) (bool, error) {
	return false, nil
}
func (m *mockRoundRepoForScore) Create(round *models.Round) error { return nil }
func (m *mockRoundRepoForScore) Update(round *models.Round) error { return nil }
func (m *mockRoundRepoForScore) Delete(id uint) error             { return nil }

type mockPlayerRepoForScore struct {
	findByID func(id uint) (*models.Player, error)
}

func (m *mockPlayerRepoForScore) FindAllByGameID(gameID uint) ([]models.Player, error) {
	return nil, nil
}
func (m *mockPlayerRepoForScore) FindByID(id uint) (*models.Player, error) { return m.findByID(id) }
func (m *mockPlayerRepoForScore) ExistsByNameAndGameID(name string, gameID uint) (bool, error) {
	return false, nil
}
func (m *mockPlayerRepoForScore) CountByGameID(gameID uint) (int64, error) { return 0, nil }
func (m *mockPlayerRepoForScore) Create(player *models.Player) error       { return nil }
func (m *mockPlayerRepoForScore) Delete(id uint) error                     { return nil }

var scoreRoundExists = &mockRoundRepoForScore{
	findByID: func(id uint) (*models.Round, error) { return &models.Round{GameID: 1}, nil },
}
var scoreRoundNotFound = &mockRoundRepoForScore{
	findByID: func(id uint) (*models.Round, error) { return nil, errors.New("record not found") },
}

func TestScore_GetAllByRound_Success(t *testing.T) {
	expected := []models.RoundScore{{Score: 10}}
	svc := NewRoundScoreService(&mockRoundScoreRepository{
		findAllByRoundID: func(roundID uint) ([]models.RoundScore, error) { return expected, nil },
	}, scoreRoundExists, &mockPlayerRepoForScore{})

	scores, err := svc.GetAllByRound(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, scores)
}

func TestScore_GetAllByRound_RoundNotFound(t *testing.T) {
	svc := NewRoundScoreService(&mockRoundScoreRepository{}, scoreRoundNotFound, &mockPlayerRepoForScore{})

	scores, err := svc.GetAllByRound(99)

	assert.ErrorIs(t, err, ErrRoundNotFound)
	assert.Nil(t, scores)
}

func TestScore_AddScore_Success(t *testing.T) {
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{
			existsByPlayerAndRoundID: func(playerID uint, roundID uint) (bool, error) { return false, nil },
			create:                   func(s *models.RoundScore) error { return nil },
		},
		scoreRoundExists,
		&mockPlayerRepoForScore{
			findByID: func(id uint) (*models.Player, error) { return &models.Player{GameID: 1}, nil },
		},
	)

	created, err := svc.AddScore(1, &models.RoundScore{PlayerID: 1, Score: 10})

	assert.NoError(t, err)
	assert.Equal(t, uint(1), created.RoundID)
}

func TestScore_AddScore_RoundNotFound(t *testing.T) {
	svc := NewRoundScoreService(&mockRoundScoreRepository{}, scoreRoundNotFound, &mockPlayerRepoForScore{})

	created, err := svc.AddScore(99, &models.RoundScore{PlayerID: 1})

	assert.ErrorIs(t, err, ErrRoundNotFound)
	assert.Nil(t, created)
}

func TestScore_AddScore_PlayerNotInGame(t *testing.T) {
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{},
		scoreRoundExists, // round.GameID = 1
		&mockPlayerRepoForScore{
			findByID: func(id uint) (*models.Player, error) { return &models.Player{GameID: 2}, nil }, // different game
		},
	)

	created, err := svc.AddScore(1, &models.RoundScore{PlayerID: 1})

	assert.ErrorIs(t, err, ErrPlayerNotInGame)
	assert.Nil(t, created)
}

func TestScore_AddScore_AlreadyExists(t *testing.T) {
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{
			existsByPlayerAndRoundID: func(playerID uint, roundID uint) (bool, error) { return true, nil },
		},
		scoreRoundExists,
		&mockPlayerRepoForScore{
			findByID: func(id uint) (*models.Player, error) { return &models.Player{GameID: 1}, nil },
		},
	)

	created, err := svc.AddScore(1, &models.RoundScore{PlayerID: 1})

	assert.ErrorIs(t, err, ErrScoreAlreadyExists)
	assert.Nil(t, created)
}

func TestScore_UpdateScore_Success(t *testing.T) {
	existing := &models.RoundScore{Score: 5}
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{
			findByID: func(id uint) (*models.RoundScore, error) { return existing, nil },
			update:   func(s *models.RoundScore) error { return nil },
		},
		scoreRoundExists,
		&mockPlayerRepoForScore{},
	)

	updated, err := svc.UpdateScore(1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 20, updated.Score)
}

func TestScore_UpdateScore_NotFound(t *testing.T) {
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{
			findByID: func(id uint) (*models.RoundScore, error) { return nil, errors.New("record not found") },
		},
		scoreRoundExists,
		&mockPlayerRepoForScore{},
	)

	updated, err := svc.UpdateScore(99, 20)

	assert.ErrorIs(t, err, ErrRoundScoreNotFound)
	assert.Nil(t, updated)
}

func TestScore_DeleteScore_Success(t *testing.T) {
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{
			findByID: func(id uint) (*models.RoundScore, error) { return &models.RoundScore{}, nil },
			delete:   func(id uint) error { return nil },
		},
		scoreRoundExists,
		&mockPlayerRepoForScore{},
	)

	assert.NoError(t, svc.DeleteScore(1))
}

func TestScore_DeleteScore_NotFound(t *testing.T) {
	svc := NewRoundScoreService(
		&mockRoundScoreRepository{
			findByID: func(id uint) (*models.RoundScore, error) { return nil, errors.New("record not found") },
		},
		scoreRoundExists,
		&mockPlayerRepoForScore{},
	)

	assert.ErrorIs(t, svc.DeleteScore(99), ErrRoundScoreNotFound)
}
