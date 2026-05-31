package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Logger is a Fiber middleware that logs each request using zerolog.
func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		statusCode := c.Response().StatusCode()
		logger := log.Info()
		if statusCode >= 500 {
			logger = log.Error()
		} else if statusCode >= 400 {
			logger = log.Warn()
		}

		logger.
			Str("method", c.Method()).
			Str("path", c.Path()).
			Int("status", statusCode).
			Dur("latency", duration).
			Str("ip", c.IP()).
			Msg("request")

		return err
	}
}
