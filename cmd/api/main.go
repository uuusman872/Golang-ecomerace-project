package main

import (
	"github.com/gin-gonic/gin"
	"github.com/uuusman872/learning-go-shop/internal/config"
	"github.com/uuusman872/learning-go-shop/internal/database"
	"github.com/uuusman872/learning-go-shop/internal/logger"
)

func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := database.New(cfg.Database)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the database")
	}
	mainDb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer mainDb.Close()
	gin.SetMode(cfg.Server.GinMode)
	log.Info().Msg("starting server")

}
