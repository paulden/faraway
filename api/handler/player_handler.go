package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/models"
	"github.com/paulden/faraway/service"
)

type PlayerHandler struct {
	service service.PlayerService
}

func NewPlayerHandler(service service.PlayerService) *PlayerHandler {
	return &PlayerHandler{service: service}
}

func (h *PlayerHandler) GetAll(c *gin.Context) {
	gameID, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	players, err := h.service.GetAllByGame(gameID)
	if err != nil {
		if errors.Is(err, service.ErrGameNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) Add(c *gin.Context) {
	gameID, err := parseID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid game id"})
		return
	}

	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.service.AddPlayer(gameID, &player)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrGameNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, service.ErrDuplicatePlayerName), errors.Is(err, service.ErrMaxPlayersReached):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusCreated, created)
}

func (h *PlayerHandler) Remove(c *gin.Context) {
	id, err := parsePlayerID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}

	if err := h.service.RemovePlayer(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func parsePlayerID(c *gin.Context) (uint, error) {
	return parseUintParam(c, "player_id")
}
