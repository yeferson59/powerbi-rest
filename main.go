package main

import (
	"context"
	"log"

	"github.com/labstack/echo/v5"
	"github.com/yeferson59/powerbi-rest/internal/config"
	"github.com/yeferson59/powerbi-rest/internal/database"
	"github.com/yeferson59/powerbi-rest/internal/handlers"
	"github.com/yeferson59/powerbi-rest/internal/middleware"
	"github.com/yeferson59/powerbi-rest/internal/routes"
)

func main() {
	ctx := context.Background()

	cfg := config.New()
	if cfg.Load() != nil {
		log.Fatal("failed to load config")
	}

	log.Println("config loaded successfully")

	db := database.NewPostgresDB(cfg.DatabaseURL).Connect(ctx)
	if db == nil {
		log.Fatal("failed to create database")
	}
	defer db.Close()

	log.Println("database created successfully")

	e, handler, middleware := echo.New(), handlers.New(db), middleware.New(db)

	if err := routes.New(e, handler, middleware).Init(); err != nil {
		e.Logger.Error("failed to initialize routes", "error", err)
	}

	if err := e.Start(":8080"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
