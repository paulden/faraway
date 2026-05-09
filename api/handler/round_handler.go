package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/service"
)

type RoundHandler struct {
	service service.RoundService
}

func NewRoundHandler(service service.RoundService) *RoundHandler {
	return &RoundHandler{service: service}
}

func (h *RoundHandler) GetAll(c *gin.Context) {
	gameID, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	rounds, err := h.service.GetAllByGame(gameID)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rounds)
}

func (h *RoundHandler) GetByID(c *gin.Context) {
	id, err := parseUintParam(c, "round_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
		return
	}

	round, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, round)
}

func (h *RoundHandler) Create(c *gin.Context) {
	gameID, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	var round models.Round
	if err := c.ShouldBindJSON(&round); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.Create(gameID, &round)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGameNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrRoundNumberExists):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *RoundHandler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "round_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
		return
	}

	var input models.Round
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.service.Update(id, &input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrRoundNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrRoundNumberExists):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *RoundHandler) Delete(c *gin.Context) {
	id, err := parseUintParam(c, "round_id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid round id"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
