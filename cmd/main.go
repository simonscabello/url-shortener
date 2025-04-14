package main

import (
	"os"
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("REDIS_URL não está definida no ambiente!")
	}

	store := storage.NewRedisStorage()

	app.Post("/encurtar", handlers.ShortenHandler(store))
	app.Get("/:slug", handlers.RedirectHandler(store))

	app.Listen(":8080")
}
