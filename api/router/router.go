package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/handler"
)

func New(gameHandler *handler.GameHandler, playerHandler *handler.PlayerHandler) *gin.Engine {
	r := gin.Default()
	registerRoutes(r, gameHandler, playerHandler)
	return r
}

func registerRoutes(r *gin.Engine, gameHandler *handler.GameHandler, playerHandler *handler.PlayerHandler) {
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	if gameHandler == nil {
		return
	}

	games := r.Group("/games")
	{
		games.GET("", gameHandler.GetAll)
		games.GET("/:id", gameHandler.GetByID)
		games.POST("", gameHandler.Create)
		games.PUT("/:id", gameHandler.Update)
		games.DELETE("/:id", gameHandler.Delete)

		games.GET("/:id/players", playerHandler.GetAll)
		games.POST("/:id/players", playerHandler.Add)
		games.DELETE("/:id/players/:player_id", playerHandler.Remove)
	}
}
