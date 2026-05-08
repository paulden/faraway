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
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type mockGameService struct {
	getAll  func() ([]models.Game, error)
	getByID func(id uint) (*models.Game, error)
	create  func(game *models.Game) error
	update  func(id uint, input *models.Game) (*models.Game, error)
	delete  func(id uint) error
}

func (m *mockGameService) GetAll() ([]models.Game, error)                        { return m.getAll() }
func (m *mockGameService) GetByID(id uint) (*models.Game, error)                  { return m.getByID(id) }
func (m *mockGameService) Create(game *models.Game) error                         { return m.create(game) }
func (m *mockGameService) Update(id uint, input *models.Game) (*models.Game, error) {
	return m.update(id, input)
}
func (m *mockGameService) Delete(id uint) error { return m.delete(id) }

func setupRouter(h *GameHandler) *gin.Engine {
	r := gin.New()
	r.GET("/games", h.GetAll)
	r.GET("/games/:id", h.GetByID)
	r.POST("/games", h.Create)
	r.PUT("/games/:id", h.Update)
	r.DELETE("/games/:id", h.Delete)
	return r
}

func TestHandler_GetAll(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		getAll: func() ([]models.Game, error) { return []models.Game{{Title: "Zelda"}}, nil },
	})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetAll_Error(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		getAll: func() ([]models.Game, error) { return nil, errors.New("db error") },
	})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games", nil))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHandler_GetByID_Found(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		getByID: func(id uint) (*models.Game, error) { return &models.Game{Title: "Zelda"}, nil },
	})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/1", nil))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_GetByID_NotFound(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		getByID: func(id uint) (*models.Game, error) { return nil, errors.New("record not found") },
	})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/99", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_GetByID_InvalidID(t *testing.T) {
	h := NewGameHandler(&mockGameService{})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/games/abc", nil))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Create(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		create: func(game *models.Game) error { return nil },
	})

	body, _ := json.Marshal(models.Game{Title: "Zelda"})
	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestHandler_Create_InvalidBody(t *testing.T) {
	h := NewGameHandler(&mockGameService{})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/games", bytes.NewBufferString("not-json")))

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_Update_Success(t *testing.T) {
	input := models.Game{Title: "Zelda: TotK", IsFinished: true}
	h := NewGameHandler(&mockGameService{
		update: func(id uint, i *models.Game) (*models.Game, error) { return &input, nil },
	})

	body, _ := json.Marshal(input)
	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/1", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestHandler_Update_NotFound(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		update: func(id uint, i *models.Game) (*models.Game, error) {
			return nil, errors.New("record not found")
		},
	})

	body, _ := json.Marshal(models.Game{Title: "Ghost"})
	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodPut, "/games/99", bytes.NewBuffer(body)))

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandler_Delete_Success(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		delete: func(id uint) error { return nil },
	})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/1", nil))

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestHandler_Delete_NotFound(t *testing.T) {
	h := NewGameHandler(&mockGameService{
		delete: func(id uint) error { return errors.New("record not found") },
	})

	w := httptest.NewRecorder()
	setupRouter(h).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/games/99", nil))

	assert.Equal(t, http.StatusNotFound, w.Code)
}
