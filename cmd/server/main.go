package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	redisCache "github.com/go-redis/cache/v9"
	"github.com/url-shortner/config"
	"github.com/url-shortner/connection"
	"github.com/url-shortner/internal/handler"
	"github.com/url-shortner/internal/repository"
	"github.com/url-shortner/internal/router"
	"github.com/url-shortner/internal/router/middleware"
	"github.com/url-shortner/internal/service/url_service"
)

func main() {
	log.SetOutput(os.Stdout)
	log.Println("URL Shortener starting...")
	appConf := config.GetAppConfigFromEnv()

	postgresConnection := connection.NewPostgresConnection(appConf.DBConfig)
	db := postgresConnection.CreateDB()

	redisConnection := connection.NewRedisConnection(appConf.RedisConfig)
	redis := redisConnection.CreateClient()
	redisCache := redisCache.New(&redisCache.Options{
		Redis: redis,
	})

	mux := chi.NewRouter()
	urlMappingRepo := repository.NewUrlMappingRepo(db)
	urlMappingService := url_service.NewUrlService(urlMappingRepo, redisCache)
	urlHandler := handler.NewURLHandler(urlMappingService)
	rateLimiter := middleware.NewRateLimiter(redis, 100, time.Minute)
	routes := router.HandlerFuncs{
		Shorten:    middleware.MakeHandler(urlHandler.Shorten),
		GetByShort: middleware.MakeHandler(urlHandler.GetByShortCode),
		URLList:    middleware.MakeHandler(urlHandler.List),
	}

	mux.Group(func(r chi.Router) {
		api := router.Create(routes, *rateLimiter)
		r.Mount("/", api)
	})

	srv := &http.Server{Addr: ":8080", Handler: mux}

	log.Println("starting server on 8080")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Server listening on :8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
