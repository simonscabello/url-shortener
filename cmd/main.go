package main

import (
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	store := storage.NewRedisStorage()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "URL Shortener API",
			"version": "1.0.0",
		})
	})

	app.Post("/encurtar", handlers.ShortenHandler(store))
	app.Get("/:slug", handlers.RedirectHandler(store))

	app.Listen(":8080")
}
