package service

import (
	"errors"
	"testing"

	"github.com/paulden/faraway/models"
	"github.com/stretchr/testify/assert"
)

type mockRoundRepository struct {
	findAllByGameID        func(gameID uint) ([]models.Round, error)
	findByID               func(id uint) (*models.Round, error)
	existsByNumberAndGameID func(number uint, gameID uint) (bool, error)
	create                 func(round *models.Round) error
	update                 func(round *models.Round) error
	delete                 func(id uint) error
}

func (m *mockRoundRepository) FindAllByGameID(gameID uint) ([]models.Round, error) {
	return m.findAllByGameID(gameID)
}
func (m *mockRoundRepository) FindByID(id uint) (*models.Round, error) { return m.findByID(id) }
func (m *mockRoundRepository) ExistsByNumberAndGameID(number uint, gameID uint) (bool, error) {
	return m.existsByNumberAndGameID(number, gameID)
}
func (m *mockRoundRepository) Create(round *models.Round) error { return m.create(round) }
func (m *mockRoundRepository) Update(round *models.Round) error { return m.update(round) }
func (m *mockRoundRepository) Delete(id uint) error             { return m.delete(id) }

type mockGameRepoForRound struct {
	findByID func(id uint) (*models.Game, error)
}

func (m *mockGameRepoForRound) FindAll() ([]models.Game, error)        { return nil, nil }
func (m *mockGameRepoForRound) FindByID(id uint) (*models.Game, error) { return m.findByID(id) }
func (m *mockGameRepoForRound) Create(game *models.Game) error         { return nil }
func (m *mockGameRepoForRound) Update(game *models.Game) error         { return nil }
func (m *mockGameRepoForRound) Delete(id uint) error                   { return nil }

var roundGameExists = &mockGameRepoForRound{
	findByID: func(id uint) (*models.Game, error) { return &models.Game{}, nil },
}
var roundGameNotFound = &mockGameRepoForRound{
	findByID: func(id uint) (*models.Game, error) { return nil, errors.New("record not found") },
}

func TestRound_GetAllByGame_Success(t *testing.T) {
	expected := []models.Round{{Number: 1}}
	svc := NewRoundService(&mockRoundRepository{
		findAllByGameID: func(gameID uint) ([]models.Round, error) { return expected, nil },
	}, roundGameExists)

	rounds, err := svc.GetAllByGame(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, rounds)
}

func TestRound_GetAllByGame_GameNotFound(t *testing.T) {
	svc := NewRoundService(&mockRoundRepository{}, roundGameNotFound)

	rounds, err := svc.GetAllByGame(99)

	assert.ErrorIs(t, err, ErrGameNotFound)
	assert.Nil(t, rounds)
}

func TestRound_GetByID_Found(t *testing.T) {
	expected := &models.Round{Number: 1}
	svc := NewRoundService(&mockRoundRepository{
		findByID: func(id uint) (*models.Round, error) { return expected, nil },
	}, roundGameExists)

	round, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, round)
}

func TestRound_GetByID_NotFound(t *testing.T) {
	svc := NewRoundService(&mockRoundRepository{
		findByID: func(id uint) (*models.Round, error) { return nil, errors.New("record not found") },
	}, roundGameExists)

	round, err := svc.GetByID(99)

	assert.ErrorIs(t, err, ErrRoundNotFound)
	assert.Nil(t, round)
}

func TestRound_Create_Success(t *testing.T) {
	round := &models.Round{Number: 1}
	svc := NewRoundService(&mockRoundRepository{
		existsByNumberAndGameID: func(number uint, gameID uint) (bool, error) { return false, nil },
		create:                  func(r *models.Round) error { return nil },
	}, roundGameExists)

	created, err := svc.Create(1, round)

	assert.NoError(t, err)
	assert.Equal(t, uint(1), created.GameID)
}

func TestRound_Create_GameNotFound(t *testing.T) {
	svc := NewRoundService(&mockRoundRepository{}, roundGameNotFound)

	created, err := svc.Create(99, &models.Round{Number: 1})

	assert.ErrorIs(t, err, ErrGameNotFound)
	assert.Nil(t, created)
}

func TestRound_Create_DuplicateNumber(t *testing.T) {
	svc := NewRoundService(&mockRoundRepository{
		existsByNumberAndGameID: func(number uint, gameID uint) (bool, error) { return true, nil },
	}, roundGameExists)

	created, err := svc.Create(1, &models.Round{Number: 1})

	assert.ErrorIs(t, err, ErrRoundNumberExists)
	assert.Nil(t, created)
}

func TestRound_Update_Success(t *testing.T) {
	existing := &models.Round{Number: 1, GameID: 1}
	svc := NewRoundService(&mockRoundRepository{
		findByID:                func(id uint) (*models.Round, error) { return existing, nil },
		existsByNumberAndGameID: func(number uint, gameID uint) (bool, error) { return false, nil },
		update:                  func(r *models.Round) error { return nil },
	}, roundGameExists)

	updated, err := svc.Update(1, &models.Round{Number: 2})

	assert.NoError(t, err)
	assert.Equal(t, uint(2), updated.Number)
}

func TestRound_Update_DuplicateNumber(t *testing.T) {
	existing := &models.Round{Number: 1, GameID: 1}
	svc := NewRoundService(&mockRoundRepository{
		findByID:                func(id uint) (*models.Round, error) { return existing, nil },
		existsByNumberAndGameID: func(number uint, gameID uint) (bool, error) { return true, nil },
	}, roundGameExists)

	updated, err := svc.Update(1, &models.Round{Number: 2})

	assert.ErrorIs(t, err, ErrRoundNumberExists)
	assert.Nil(t, updated)
}

func TestRound_Delete_Success(t *testing.T) {
	svc := NewRoundService(&mockRoundRepository{
		findByID: func(id uint) (*models.Round, error) { return &models.Round{}, nil },
		delete:   func(id uint) error { return nil },
	}, roundGameExists)

	assert.NoError(t, svc.Delete(1))
}

func TestRound_Delete_NotFound(t *testing.T) {
	svc := NewRoundService(&mockRoundRepository{
		findByID: func(id uint) (*models.Round, error) { return nil, errors.New("record not found") },
	}, roundGameExists)

	assert.ErrorIs(t, svc.Delete(99), ErrRoundNotFound)
}
