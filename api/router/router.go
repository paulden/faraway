package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/handler"
)

func New(
	gameHandler *handler.GameHandler,
	playerHandler *handler.PlayerHandler,
	roundHandler *handler.RoundHandler,
	scoreHandler *handler.RoundScoreHandler,
) *gin.Engine {
	r := gin.Default()
	registerRoutes(r, gameHandler, playerHandler, roundHandler, scoreHandler)
	return r
}

func registerRoutes(
	r *gin.Engine,
	gameHandler *handler.GameHandler,
	playerHandler *handler.PlayerHandler,
	roundHandler *handler.RoundHandler,
	scoreHandler *handler.RoundScoreHandler,
) {
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

		games.GET("/:id/rounds", roundHandler.GetAll)
		games.POST("/:id/rounds", roundHandler.Create)
		games.GET("/:id/rounds/:round_id", roundHandler.GetByID)
		games.PUT("/:id/rounds/:round_id", roundHandler.Update)
		games.DELETE("/:id/rounds/:round_id", roundHandler.Delete)

		games.GET("/:id/rounds/:round_id/scores", scoreHandler.GetAll)
		games.POST("/:id/rounds/:round_id/scores", scoreHandler.Add)
		games.PUT("/:id/rounds/:round_id/scores/:score_id", scoreHandler.Update)
		games.DELETE("/:id/rounds/:round_id/scores/:score_id", scoreHandler.Delete)
	}
}
