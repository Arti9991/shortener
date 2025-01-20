package auth

import (
	"context"
	"net/http"
)

type KeyContext string

var key = KeyContext("UserID")

func MiddlewareAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		UserID := "125"
		ctx := context.WithValue(req.Context(), key, UserID)
		req = req.WithContext(ctx)
		// передаём управление хендлеру
		h.ServeHTTP(res, req)
	}
}
