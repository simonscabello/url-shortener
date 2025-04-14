package main

import (
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	store := storage.NewRedisStorage()

	app.Post("/encurtar", handlers.ShortenHandler(store))
	app.Get("/:slug", handlers.RedirectHandler(store))

	app.Listen(":8080")
}
