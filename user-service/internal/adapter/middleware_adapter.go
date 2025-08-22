package adapter

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
	"user-service/config"
	"user-service/internal/adapter/handler/response"
)

type MiddlewareAdapterInterface interface {
	CheckToken() echo.MiddlewareFunc
}

type middlewareAdapter struct {
	cfg *config.Config
}

func NewMiddlewareAdapter(cfg *config.Config) MiddlewareAdapterInterface {
	return &middlewareAdapter{
		cfg: cfg,
	}
}

func (m *middlewareAdapter) CheckToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			respError := response.DefaultResponse{}
			redisConn := config.NewRedisClient()
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				respError.Message = "missing or invalid token"
				respError.Data = nil
				return c.JSON(http.StatusUnauthorized, respError)
			}
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			getSession, err := redisConn.HGetAll(c.Request().Context(), tokenString).Result()
			if err != nil || len(getSession) == 0 {
				respError.Message = err.Error()
				respError.Data = nil
				return c.JSON(http.StatusUnauthorized, respError)
			}

			c.Set("user", getSession)
			return next(c)
		}
	}
}
