package handlers

import (
	"time"
	"url-shortener/internal/storage"
	"url-shortener/internal/utils"

	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type request struct {
	URL string `json:"url"`
	TTL int    `json:"ttl"` // TTL opcional em segundos
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

var validate = validator.New()

func ShortenHandler(store storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body request

		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Erro ao ler corpo da requisição",
			})
		}

		if !strings.HasPrefix(body.URL, "http://") && !strings.HasPrefix(body.URL, "https://") {
			body.URL = "https://" + body.URL
		}

		if err := validate.Struct(body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "URL inválida",
			})
		}

		var slug string
		for {
			slug = utils.GenerateSlug(6)
			if _, exists := store.Get(slug); !exists {
				break
			}
		}

		var ttl time.Duration
		if body.TTL > 0 {
			ttl = time.Duration(body.TTL) * time.Second
		}

		store.Save(slug, body.URL, ttl)

		return c.JSON(ShortenResponse{
			ShortURL: c.BaseURL() + "/" + slug,
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

		store.IncrementClicks(slug)

		return c.Redirect(url, fiber.StatusMovedPermanently)
	}
}

func StatsHandler(store storage.Storage) fiber.Handler {
	return func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		clicks := store.GetClicks(slug)

		return c.JSON(fiber.Map{
			"slug":   slug,
			"clicks": clicks,
		})
	}
}
