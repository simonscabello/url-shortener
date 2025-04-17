package handlers_test

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"url-shortener/internal/handlers"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type MockStorage struct {
	data   map[string]string
	clicks map[string]int
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		data:   make(map[string]string),
		clicks: make(map[string]int),
	}
}

func (m *MockStorage) Save(slug, url string, ttl time.Duration) {
	m.data[slug] = url
}

func (m *MockStorage) Get(slug string) (string, bool) {
	val, ok := m.data[slug]
	return val, ok
}

func (m *MockStorage) IncrementClicks(slug string) {
	m.clicks[slug]++
}

func (m *MockStorage) GetClicks(slug string) int {
	return m.clicks[slug]
}

func TestShortenHandler(t *testing.T) {
	app := fiber.New()
	store := NewMockStorage()
	app.Post("/", handlers.ShortenHandler(store))

	body := `{"url": "example.com", "ttl": 3600}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestRedirectHandler(t *testing.T) {
	app := fiber.New()
	store := NewMockStorage()
	store.Save("test", "https://example.com", 0)
	app.Get("/:slug", handlers.RedirectHandler(store))

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 301, resp.StatusCode)
}

func TestStatsHandler(t *testing.T) {
	app := fiber.New()
	store := NewMockStorage()
	store.clicks["abc"] = 5
	app.Get("/stats/:slug", handlers.StatsHandler(store))

	req := httptest.NewRequest("GET", "/stats/abc", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
