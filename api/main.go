package main

import (
	"log"

	"github.com/paulden/faraway/config"
	"github.com/paulden/faraway/db"
	"github.com/paulden/faraway/handler"
	"github.com/paulden/faraway/repository"
	"github.com/paulden/faraway/router"
	"github.com/paulden/faraway/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	gameRepo := repository.NewGameRepository(database)
	gameService := service.NewGameService(gameRepo)
	gameHandler := handler.NewGameHandler(gameService)

	playerRepo := repository.NewPlayerRepository(database)
	playerService := service.NewPlayerService(playerRepo, gameRepo)
	playerHandler := handler.NewPlayerHandler(playerService)

	roundRepo := repository.NewRoundRepository(database)
	roundService := service.NewRoundService(roundRepo, gameRepo)
	roundHandler := handler.NewRoundHandler(roundService)

	scoreRepo := repository.NewRoundScoreRepository(database)
	scoreService := service.NewRoundScoreService(scoreRepo, roundRepo, playerRepo)
	scoreHandler := handler.NewRoundScoreHandler(scoreService)

	r := router.New(gameHandler, playerHandler, roundHandler, scoreHandler)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
