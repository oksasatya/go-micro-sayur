package app

import (
	"context"
	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/config"
	"user-service/internal/adapter/handler"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/service"
	"user-service/utils/validator"
)

func RunServer() {
	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres()
	if err != nil {
		log.Fatalf("[RunServer-1] RunServer: failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db.DB)
	jwtService := service.NewJwtService(cfg)
	userService := service.NewUserService(userRepo, cfg, jwtService)

	e := echo.New()
	e.Use(middleware.CORS())

	customValidator := validator.NewValidator()
	err = en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)
	if err != nil {
		log.Fatalf("[RunServer-2] RunServer: failed to register translations: %v", err)
	}
	e.Validator = customValidator
	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})
	handler.NewUserHandler(e, userService, cfg)
	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}
		err = e.Start(":" + cfg.App.AppPort)
		if err != nil {
			log.Fatalf("[RunServer-3] RunServer: failed to start server: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	log.Printf("[RunServer-4] RunServer: shutting down server gracefully...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = e.Shutdown(ctx)
	if err != nil {
		log.Fatalf("[RunServer-5] RunServer: failed to shutdown server: %v", err)
	} else {
		log.Printf("[RunServer-6] RunServer: server shutdown successfully")
	}
}
