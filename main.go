package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/yeferson59/powerbi-rest/internal/config"
	"github.com/yeferson59/powerbi-rest/internal/database"
	"github.com/yeferson59/powerbi-rest/internal/handlers"
	"github.com/yeferson59/powerbi-rest/internal/metrics"
	"github.com/yeferson59/powerbi-rest/internal/middleware"
	"github.com/yeferson59/powerbi-rest/internal/routes"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	if err := cfg.Load(); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	log.Println("config loaded successfully")

	db, err := database.NewPostgresDB(cfg.DatabaseURL).Connect(ctx)
	if err != nil {
		log.Fatalf("failed to create database: %v", err)
	}
	defer db.Close()

	log.Println("database created successfully")

	metricsStore := metrics.NewStore(db)
	e := echo.New()
	handler := handlers.New(metricsStore)
	mw := middleware.New(metricsStore)

	if err := routes.New(e, handler, mw).Init(); err != nil {
		e.Logger.Error("failed to initialize routes", "error", err)
		return
	}

	if err := e.Start(":" + cfg.Port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
