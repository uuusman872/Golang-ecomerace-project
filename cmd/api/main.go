package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/uuusman872/learning-go-shop/docs"
	"github.com/uuusman872/learning-go-shop/internal/config"
	"github.com/uuusman872/learning-go-shop/internal/database"
	"github.com/uuusman872/learning-go-shop/internal/events"
	"github.com/uuusman872/learning-go-shop/internal/interfaces"
	"github.com/uuusman872/learning-go-shop/internal/logger"
	"github.com/uuusman872/learning-go-shop/internal/providers"
	"github.com/uuusman872/learning-go-shop/internal/server"
	service "github.com/uuusman872/learning-go-shop/internal/services"
)

// @title           Learning Go Shop API
// @version         1.0
// @description     Ecommerce shop API built with Gin.
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @security BearerAuth
func main() {
	log := logger.New()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}
	db, err := database.New(&cfg.Database)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the database")
	}
	mainDb, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database connection")
	}
	defer func() {
		err := mainDb.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to close the mainDB")
		}
	}()

	gin.SetMode(cfg.Server.GinMode)

	ctx := context.Background()
	eventPublisher, err := events.NewEventPublisher(ctx, &cfg.AWS)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize the event publisher")
		return
	}

	authService := service.NewAuthService(db, cfg, eventPublisher)
	productService := service.NewProductService(db)
	userService := service.NewUserService(db)
	cartService := service.NewCartService(db)
	orderService := service.NewOrderService(db)
	var uploadProvider interfaces.UploadProvider

	if cfg.Upload.UploadProvider == "s3" {
		uploadProvider = providers.NewS3Provider(cfg)
	} else {
		uploadProvider = providers.NewLocalUploadProvider(cfg.Upload.Path)
	}
	uploadService := service.NewUploadService(uploadProvider)
	srv := server.New(cfg, db, &log, authService, productService, userService, uploadService, cartService, orderService)
	router := srv.SetupRoutes()

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Info().Str("port", cfg.Server.Port).Msg("starting http server")
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("failed to start http server")
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("The server is shutting down")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to shutdown http server")
	}
	log.Info().Msg("shutting down database")
}
