package handlers

import (
	"time"
	"url-shortener/internal/storage"
	"url-shortener/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"` // TTL opcional em segundos
}

func ShortenHandler(store storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body request
		if err := c.BodyParser(&body); err != nil || body.URL == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "URL inválida",
			})
		}

		slug := utils.GenerateSlug(6)

		var ttl time.Duration
		if body.TTL > 0 {
			ttl = time.Duration(body.TTL) * time.Second
		}

		store.Save(slug, body.URL, ttl)

		return c.JSON(fiber.Map{
			"short_url": c.BaseURL() + "/" + slug,
		})
	}
}

func RedirectHandler(store storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		url, found := store.Get(slug)

		if !found {
			return c.Status(fiber.StatusNotFound).SendString("URL não encontrada")
		}

		return c.Redirect(url, fiber.StatusMovedPermanently)
	}
}
