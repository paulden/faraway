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

type mockRoundService struct {
	getAllByGame func(gameID uint) ([]models.Round, error)
	getByID     func(id uint) (*models.Round, error)
	create      func(gameID uint, round *models.Round) (*models.Round, error)
	update      func(id uint, input *models.Round) (*models.Round, error)
	delete      func(id uint) error
}

func (m *mockRoundService) GetAllByGame(gameID uint) ([]models.Round, error) {
	return m.getAllByGame(gameID)
}
func (m *mockRoundService) GetByID(id uint) (*models.Round, error) { return m.getByID(id) }
func (m *mockRoundService) Create(gameID uint, round *models.Round) (*models.Round, error) {
	return m.create(gameID, round)
}
func (m *mockRoundService) Update(id uint, input *models.Round) (*models.Round, error) {
	return m.update(id, input)
}
func (m *mockRoundService) Delete(id uint) error { return m.delete(id) }

func setupRoundRouter(h *RoundHandler) *gin.Engine {
	r := gin.New()
	r.GET("/games/:id/rounds", h.GetAll)
	r.GET("/games/:id/rounds/:round_id", h.GetByID)
	r.POST("/games/:id/rounds", h.Create)
	r.PUT("/games/:id/rounds/:round_id", h.Update)
	r.DELETE("/games/:id/rounds/:round_id", h.Delete)
	return r
}

func TestRoundHandler_GetAll_Success(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		getAllByGame: func(gameID uint) ([]models.Round, error) { return []models.Round{{Number: 1}}, nil },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoundHandler_GetAll_GameNotFound(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		getAllByGame: func(gameID uint) ([]models.Round, error) { return nil, service.ErrGameNotFound },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/99/rounds", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoundHandler_GetByID_Found(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		getByID: func(id uint) (*models.Round, error) { return &models.Round{Number: 1}, nil },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds/1", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoundHandler_GetByID_NotFound(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		getByID: func(id uint) (*models.Round, error) { return nil, service.ErrRoundNotFound },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds/99", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoundHandler_Create_Success(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		create: func(gameID uint, round *models.Round) (*models.Round, error) {
			return &models.Round{Number: round.Number, GameID: gameID}, nil
		},
	})

	body, _ := json.Marshal(models.Round{Number: 1})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRoundHandler_Create_DuplicateNumber(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		create: func(gameID uint, round *models.Round) (*models.Round, error) {
			return nil, service.ErrRoundNumberExists
		},
	})

	body, _ := json.Marshal(models.Round{Number: 1})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRoundHandler_Update_Success(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		update: func(id uint, input *models.Round) (*models.Round, error) {
			return &models.Round{Number: input.Number}, nil
		},
	})

	body, _ := json.Marshal(models.Round{Number: 2})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1/rounds/1", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRoundHandler_Update_NotFound(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		update: func(id uint, input *models.Round) (*models.Round, error) {
			return nil, service.ErrRoundNotFound
		},
	})

	body, _ := json.Marshal(models.Round{Number: 2})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1/rounds/99", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoundHandler_Delete_Success(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		delete: func(id uint) error { return nil },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/rounds/1", nil))

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestRoundHandler_Delete_NotFound(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		delete: func(id uint) error { return service.ErrRoundNotFound },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/rounds/99", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoundHandler_InvalidID(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/abc/rounds", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_Create_GameNotFound(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		create: func(gameID uint, round *models.Round) (*models.Round, error) {
			return nil, service.ErrGameNotFound
		},
	})

	body, _ := json.Marshal(models.Round{Number: 1})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/99/rounds", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRoundHandler_Update_DuplicateNumber(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		update: func(id uint, input *models.Round) (*models.Round, error) {
			return nil, service.ErrRoundNumberExists
		},
	})

	body, _ := json.Marshal(models.Round{Number: 2})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1/rounds/1", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRoundHandler_Create_InvalidBody(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/1/rounds", bytes.NewBufferString("not-json")))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_Update_InvalidID(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	body, _ := json.Marshal(models.Round{Number: 2})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1/rounds/abc", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_Delete_InvalidID(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1/rounds/abc", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_GetByID_InvalidID(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds/abc", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_Create_InvalidGameID(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	body, _ := json.Marshal(models.Round{Number: 1})
	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games/abc/rounds", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_GetAll_InvalidGameID(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/abc/rounds", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRoundHandler_GetAll_Error(t *testing.T) {
	h := NewRoundHandler(&mockRoundService{
		getAllByGame: func(gameID uint) ([]models.Round, error) { return nil, errors.New("db error") },
	})

	w := httptest.NewRecorder()
	setupRoundRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1/rounds", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
