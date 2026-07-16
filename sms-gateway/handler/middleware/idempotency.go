package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

const IdempotencyKeyHeader = "X-Idempotency-Key"

func RequireIdempotencyKey(c fiber.Ctx) error {
	if fiber.IsMethodSafe(c.Method()) {
		return c.Next()
	}

	key := strings.TrimSpace(c.Get(IdempotencyKeyHeader))
	if key == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "X-Idempotency-Key header is required",
		})
	}

	return c.Next()
}
