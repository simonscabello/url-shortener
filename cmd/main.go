package main

import (
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	app := fiber.New()

	store := storage.NewRedisStorage()

	app.Get("/", handlers.HelloHandler())

	app.Post("/encurtar", handlers.ShortenHandler(store))
	app.Get("/:slug", handlers.RedirectHandler(store))
	app.Get("/:slug/stats", handlers.StatsHandler(store))

	app.Listen(":8080")
}
