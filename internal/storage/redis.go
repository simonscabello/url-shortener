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
	IncrementClicks(slug string)
	GetClicks(slug string) int
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
	err := r.client.Set(r.ctx, slug, url, ttl).Err()
	if err != nil {
		panic("Erro ao salvar slug no Redis: " + err.Error())
	}
}

func (r *RedisStorage) Get(slug string) (string, bool) {
	val, err := r.client.Get(r.ctx, slug).Result()
	if err == redis.Nil || err != nil {
		return "", false
	}
	return val, true
}

func (r *RedisStorage) IncrementClicks(slug string) {
	r.client.Incr(r.ctx, slug+":clicks")
}

func (r *RedisStorage) GetClicks(slug string) int {
	val, err := r.client.Get(r.ctx, slug+":clicks").Int()
	if err != nil {
		return 0
	}
	return val
}
