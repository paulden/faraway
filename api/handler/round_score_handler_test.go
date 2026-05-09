package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/service"
	"github.com/stretchr/testify/assert"
)

type mockRoundScoreService struct {
	getAllByRound func(roundID uint) ([]models.RoundScore, error)
	addScore     func(roundID uint, score *models.RoundScore) (*models.RoundScore, error)
	updateScore  func(id uint, newScore int) (*models.RoundScore, error)
	deleteScore  func(id uint) error
}

func (m *mockRoundScoreService) GetAllByRound(roundID uint) ([]models.RoundScore, error) {
	return m.getAllByRound(roundID)
}
func (m *mockRoundScoreService) AddScore(roundID uint, score *models.RoundScore) (*models.RoundScore, error) {
	return m.addScore(roundID, score)
}
func (m *mockRoundScoreService) UpdateScore(id uint, newScore int) (*models.RoundScore, error) {
	return m.updateScore(id, newScore)
}
func (m *mockRoundScoreService) DeleteScore(id uint) error { return m.deleteScore(id) }

func setupScoreRouter(h *RoundScoreHandler) *gin.Engine {
	r := gin.New()
	r.GET("/games/:id/rounds/:round_id/scores", h.GetAll)
	r.POST("/games/:id/rounds/:round_id/scores", h.Add)
	r.PUT("/games/:id/rounds/:round_id/scores/:score_id", h.Update)
	r.DELETE("/games/:id/rounds/:round_id/scores/:score_id", h.Delete)
	return r
}

func TestScoreHandler_GetAll_Success(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		getAllByRound: func(roundID uint) ([]models.RoundScore, error) {
			return []models.RoundScore{{Score: 10}}, nil
		},
	})

	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds/1/scores", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestScoreHandler_GetAll_RoundNotFound(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		getAllByRound: func(roundID uint) ([]models.RoundScore, error) {
			return nil, service.ErrRoundNotFound
		},
	})

	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds/99/scores", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestScoreHandler_GetAll_InvalidRoundID(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{})

	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds/abc/scores", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScoreHandler_Add_Success(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		addScore: func(roundID uint, score *models.RoundScore) (*models.RoundScore, error) {
			return &models.RoundScore{Score: score.Score, RoundID: roundID}, nil
		},
	})

	body, _ := json.Marshal(models.RoundScore{PlayerID: 1, Score: 10})
	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds/1/scores", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestScoreHandler_Add_RoundNotFound(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		addScore: func(roundID uint, score *models.RoundScore) (*models.RoundScore, error) {
			return nil, service.ErrRoundNotFound
		},
	})

	body, _ := json.Marshal(models.RoundScore{PlayerID: 1, Score: 10})
	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds/99/scores", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestScoreHandler_Add_PlayerNotInGame(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		addScore: func(roundID uint, score *models.RoundScore) (*models.RoundScore, error) {
			return nil, service.ErrPlayerNotInGame
		},
	})

	body, _ := json.Marshal(models.RoundScore{PlayerID: 1, Score: 10})
	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds/1/scores", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestScoreHandler_Add_AlreadyExists(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		addScore: func(roundID uint, score *models.RoundScore) (*models.RoundScore, error) {
			return nil, service.ErrScoreAlreadyExists
		},
	})

	body, _ := json.Marshal(models.RoundScore{PlayerID: 1, Score: 10})
	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds/1/scores", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestScoreHandler_Add_InvalidBody(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{})

	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds/1/scores", bytes.NewBufferString("not-json")))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestScoreHandler_Update_Success(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		updateScore: func(id uint, newScore int) (*models.RoundScore, error) {
			return &models.RoundScore{Score: newScore}, nil
		},
	})

	body, _ := json.Marshal(map[string]int{"score": 25})
	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1/rounds/1/scores/1", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestScoreHandler_Update_NotFound(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		updateScore: func(id uint, newScore int) (*models.RoundScore, error) {
			return nil, service.ErrRoundScoreNotFound
		},
	})

	body, _ := json.Marshal(map[string]int{"score": 25})
	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1/rounds/1/scores/99", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestScoreHandler_Delete_Success(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		deleteScore: func(id uint) error { return nil },
	})

	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/rounds/1/scores/1", nil))

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestScoreHandler_Delete_NotFound(t *testing.T) {
	h := NewRoundScoreHandler(&mockRoundScoreService{
		deleteScore: func(id uint) error { return service.ErrRoundScoreNotFound },
	})

	w := httptest.NewRecorder()
	setupScoreRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/rounds/1/scores/99", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}
