package storage

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Storage interface {
	Save(slug, url string, ttl time.Duration)
	Get(slug string) (string, bool)
}

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage() *RedisStorage {
	url := os.Getenv("REDIS_URL")
	opts, err := redis.ParseURL(url)
	if err != nil {
		panic("Erro ao parsear REDIS_URL: " + err.Error())
	}

	rdb := redis.NewClient(opts)

	return &RedisStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisStorage) Save(slug, url string, ttl time.Duration) {
	r.client.Set(r.ctx, slug, url, ttl)
}

func (r *RedisStorage) Get(slug string) (string, bool) {
	val, err := r.client.Get(r.ctx, slug).Result()
	if err == redis.Nil || err != nil {
		return "", false
	}
	return val, true
}
