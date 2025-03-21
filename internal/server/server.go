package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/mangadev-courses/goexpert-rate-limiter/internal/server/middleware"
	"github.com/mangadev-courses/goexpert-rate-limiter/pkg/limiter"
	"github.com/mangadev-courses/goexpert-rate-limiter/pkg/redis"
)

type Server struct {
	Echo *echo.Echo
}

func New(e *echo.Echo, ctx context.Context) error {
	err := setupMiddleware(e, ctx)
	if err != nil {
		return fmt.Errorf("Error setting up middleware: %w", err)
	}

	setupRouter(e)

	return nil
}

func setupMiddleware(e *echo.Echo, ctx context.Context) error {
	e.Use(echoMiddleware.Recover())

	redisClient, err := redis.NewRedisClient(ctx)
	if err != nil {
		return fmt.Errorf("Error connecting to Redis: %w", err)
	}
	limiter := limiter.New(redisClient)
	e.Use(middleware.Limiter(limiter))

	return nil
}

func setupRouter(e *echo.Echo) {
	e.GET("/healthz", func(c echo.Context) error {
		return c.NoContent(200)
	})

	e.Any("*", func(c echo.Context) error {
		return c.NoContent(http.StatusNotFound)
	})
}
