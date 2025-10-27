package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/url-shortner/internal/router/middleware"
)

type HandlerFuncs struct {
	Shorten    http.HandlerFunc
	GetByShort http.HandlerFunc
	URLList    http.HandlerFunc
}

func Create(handlerFuncs HandlerFuncs, rateLimiter middleware.RateLimiter) http.Handler {
	api := chi.NewRouter()
	api.Group(func(r chi.Router) {
		r.Route("/url", func(r chi.Router) {
			r.Use(rateLimiter.LimitByIP)
			r.Post("/", handlerFuncs.Shorten)
			r.Get("/{shortCode}", handlerFuncs.GetByShort)
		})
		r.Route("/admin", func(r chi.Router) {
			r.Get("/urls", handlerFuncs.URLList)
		})
	})
	return api
}
