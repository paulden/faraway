package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulden/faraway/handler"
)

func New(gameHandler *handler.GameHandler) *gin.Engine {
	r := gin.Default()
	registerRoutes(r, gameHandler)
	return r
}

func registerRoutes(r *gin.Engine, gameHandler *handler.GameHandler) {
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
	}
}
