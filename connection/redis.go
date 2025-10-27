package connection

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/url-shortner/config"
)

type RedisConnection interface {
	CreateClient() *redis.Client
	GetClient() *redis.Client
	Ping() error
	Close() error
}

type redisConnection struct {
	Config config.RedisConfig
	Client *redis.Client
}

// Constructor
func NewRedisConnection(cfg config.RedisConfig) RedisConnection {
	return &redisConnection{
		Config: cfg,
	}
}

// CreateClient initializes the Redis client
func (r *redisConnection) CreateClient() *redis.Client {
	addr := fmt.Sprintf("%s:%d", r.Config.Host, r.Config.Port)
	r.Client = redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           r.Config.DB,
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
		DialTimeout:  time.Second * 5,
	})

	if err := r.Ping(); err != nil {
		panic(fmt.Sprintf("failed to connect to Redis at %s: %v", addr, err))
	}
	return r.Client
}

// GetClient returns the existing client
func (r *redisConnection) GetClient() *redis.Client {
	if r.Client == nil {
		return r.CreateClient()
	}
	return r.Client
}

// Ping checks connectivity
func (r *redisConnection) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.Client.Ping(ctx).Result()
	return err
}

// Close shuts down the client
func (r *redisConnection) Close() error {
	if r.Client != nil {
		return r.Client.Close()
	}
	return nil
}
