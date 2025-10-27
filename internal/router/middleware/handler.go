package middleware

import (
	"net/http"
)

// MakeHandler Create handler
func MakeHandler(handlerFunc func() http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlerFunc()(w, r)
	}
}
