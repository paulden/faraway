package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/service"
	"github.com/stretchr/testify/assert"
)

type mockPlayerService struct {
	getAllByGame func(gameID uint) ([]models.Player, error)
	addPlayer   func(gameID uint, player *models.Player) (*models.Player, error)
	removePlayer func(id uint) error
}

func (m *mockPlayerService) GetAllByGame(gameID uint) ([]models.Player, error) {
	return m.getAllByGame(gameID)
}
func (m *mockPlayerService) AddPlayer(gameID uint, player *models.Player) (*models.Player, error) {
	return m.addPlayer(gameID, player)
}
func (m *mockPlayerService) RemovePlayer(id uint) error { return m.removePlayer(id) }

func setupPlayerRouter(h *PlayerHandler) *gin.Engine {
	r := gin.New()
	r.GET("/games/:id/players", h.GetAll)
	r.POST("/games/:id/players", h.Add)
	r.DELETE("/games/:id/players/:player_id", h.Remove)
	return r
}

func TestPlayerHandler_GetAll_Success(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		getAllByGame: func(gameID uint) ([]models.Player, error) {
			return []models.Player{{Name: "Alice"}}, nil
		},
	})

	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/players", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestPlayerHandler_GetAll_GameNotFound(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		getAllByGame: func(gameID uint) ([]models.Player, error) {
			return nil, service.ErrGameNotFound
		},
	})

	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/99/players", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlayerHandler_GetAll_InvalidGameID(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{})

	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/abc/players", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlayerHandler_Add_Success(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		addPlayer: func(gameID uint, player *models.Player) (*models.Player, error) {
			return &models.Player{Name: "Alice", GameID: gameID}, nil
		},
	})

	body, _ := json.Marshal(models.Player{Name: "Alice"})
	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/players", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestPlayerHandler_Add_GameNotFound(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		addPlayer: func(gameID uint, player *models.Player) (*models.Player, error) {
			return nil, service.ErrGameNotFound
		},
	})

	body, _ := json.Marshal(models.Player{Name: "Alice"})
	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/99/players", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestPlayerHandler_Add_DuplicateName(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		addPlayer: func(gameID uint, player *models.Player) (*models.Player, error) {
			return nil, service.ErrDuplicatePlayerName
		},
	})

	body, _ := json.Marshal(models.Player{Name: "Alice"})
	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/players", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestPlayerHandler_Add_MaxPlayersReached(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		addPlayer: func(gameID uint, player *models.Player) (*models.Player, error) {
			return nil, service.ErrMaxPlayersReached
		},
	})

	body, _ := json.Marshal(models.Player{Name: "Alice"})
	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/players", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestPlayerHandler_Add_InvalidBody(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{})

	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/players", bytes.NewBufferString("not-json")))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestPlayerHandler_Remove_Success(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		removePlayer: func(id uint) error { return nil },
	})

	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/players/2", nil))

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestPlayerHandler_Remove_NotFound(t *testing.T) {
	h := NewPlayerHandler(&mockPlayerService{
		removePlayer: func(id uint) error { return errors.New("player not found") },
	})

	w := httptest.NewRecorder()
	setupPlayerRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/players/99", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}
