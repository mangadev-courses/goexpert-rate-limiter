package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/mangadev-courses/goexpert-rate-limiter/pkg/limiter"
)

func Limiter(limiter *limiter.Limiter) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			apiKey := c.Request().Header.Get("API_KEY")

			res, err := limiter.Allow(c.Request().Context(), c.RealIP(), apiKey)
			if err != nil {
				return echo.NewHTTPError(500, err.Error())
			}

			if !res.Allowed {
				return echo.NewHTTPError(429, "you have reached the maximum number of requests or actions allowed within a certain time frame")
			}

			c.Response().Header().Set("X-RateLimit-Remaining", fmt.Sprint(res.Remaining))

			return next(c)
		}
	}
}
