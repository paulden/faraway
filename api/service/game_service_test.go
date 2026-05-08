package service

import (
	"errors"
	"testing"

	"github.com/paulden/faraway/models"
	"github.com/stretchr/testify/assert"
)

type mockGameRepository struct {
	findAll  func() ([]models.Game, error)
	findByID func(id uint) (*models.Game, error)
	create   func(game *models.Game) error
	update   func(game *models.Game) error
	delete   func(id uint) error
}

func (m *mockGameRepository) FindAll() ([]models.Game, error)       { return m.findAll() }
func (m *mockGameRepository) FindByID(id uint) (*models.Game, error) { return m.findByID(id) }
func (m *mockGameRepository) Create(game *models.Game) error         { return m.create(game) }
func (m *mockGameRepository) Update(game *models.Game) error         { return m.update(game) }
func (m *mockGameRepository) Delete(id uint) error                   { return m.delete(id) }

func TestGetAll(t *testing.T) {
	expected := []models.Game{{Title: "Zelda"}}
	svc := NewGameService(&mockGameRepository{
		findAll: func() ([]models.Game, error) { return expected, nil },
	})

	games, err := svc.GetAll()

	assert.NoError(t, err)
	assert.Equal(t, expected, games)
}

func TestGetByID_Found(t *testing.T) {
	expected := &models.Game{Title: "Zelda"}
	svc := NewGameService(&mockGameRepository{
		findByID: func(id uint) (*models.Game, error) { return expected, nil },
	})

	game, err := svc.GetByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expected, game)
}

func TestGetByID_NotFound(t *testing.T) {
	svc := NewGameService(&mockGameRepository{
		findByID: func(id uint) (*models.Game, error) { return nil, errors.New("record not found") },
	})

	game, err := svc.GetByID(99)

	assert.Error(t, err)
	assert.Nil(t, game)
}

func TestCreate(t *testing.T) {
	game := &models.Game{Title: "Zelda"}
	svc := NewGameService(&mockGameRepository{
		create: func(g *models.Game) error { return nil },
	})

	err := svc.Create(game)

	assert.NoError(t, err)
}

func TestUpdate_Success(t *testing.T) {
	existing := &models.Game{Title: "Zelda", IsFinished: false}
	svc := NewGameService(&mockGameRepository{
		findByID: func(id uint) (*models.Game, error) { return existing, nil },
		update:   func(g *models.Game) error { return nil },
	})

	updated, err := svc.Update(1, &models.Game{Title: "Zelda: TotK", IsFinished: true})

	assert.NoError(t, err)
	assert.Equal(t, "Zelda: TotK", updated.Title)
	assert.True(t, updated.IsFinished)
}

func TestUpdate_NotFound(t *testing.T) {
	svc := NewGameService(&mockGameRepository{
		findByID: func(id uint) (*models.Game, error) { return nil, errors.New("record not found") },
	})

	updated, err := svc.Update(99, &models.Game{Title: "Ghost"})

	assert.Error(t, err)
	assert.Nil(t, updated)
}

func TestDelete_Success(t *testing.T) {
	svc := NewGameService(&mockGameRepository{
		findByID: func(id uint) (*models.Game, error) { return &models.Game{}, nil },
		delete:   func(id uint) error { return nil },
	})

	err := svc.Delete(1)

	assert.NoError(t, err)
}

func TestDelete_NotFound(t *testing.T) {
	svc := NewGameService(&mockGameRepository{
		findByID: func(id uint) (*models.Game, error) { return nil, errors.New("record not found") },
	})

	err := svc.Delete(99)

	assert.Error(t, err)
}
