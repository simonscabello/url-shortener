package main

import (
	"fmt"
	"os"
	"url-shortener/internal/handlers"
	"url-shortener/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("🔍 Variáveis de ambiente disponíveis:")
	for _, e := range os.Environ() {
		fmt.Println(e)
	}

	app := fiber.New()

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		panic("REDIS_URL não está definida no ambiente!")
	}

	fmt.Println("REDIS_URL:", redisURL)

	store := storage.NewRedisStorage()

	app.Post("/encurtar", handlers.ShortenHandler(store))
	app.Get("/:slug", handlers.RedirectHandler(store))

	app.Listen(":8080")
}
