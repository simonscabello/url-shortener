package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func HelloHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "URL Shortener API",
			"version": "1.0.0",
			"author":  "Simon",
		})
	}
}
