package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/service"
)

type RoundScoreHandler struct {
	service service.RoundScoreService
}

func NewRoundScoreHandler(service service.RoundScoreService) *RoundScoreHandler {
	return &RoundScoreHandler{service: service}
}

func (h *RoundScoreHandler) GetAll(c *gin.Context) {
	roundID, err := parseUintParam(c, "round_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
		return
	}

	scores, err := h.service.GetAllByRound(roundID)
	if err != nil {
		if errors.Is(err, service.ErrRoundNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, scores)
}

func (h *RoundScoreHandler) Add(c *gin.Context) {
	roundID, err := parseUintParam(c, "round_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
		return
	}

	var score models.RoundScore
	if err := c.ShouldBindJSON(&score); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.AddScore(roundID, &score)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoundNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrPlayerNotInGame), errors.Is(err, service.ErrScoreAlreadyExists):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *RoundScoreHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "score_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid score id"})
		return
	}

	var body struct {
		Score int `json:"score"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.UpdateScore(id, body.Score)
	if err != nil {
		if errors.Is(err, service.ErrRoundScoreNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *RoundScoreHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "score_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid score id"})
		return
	}

	if err := h.service.DeleteScore(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
