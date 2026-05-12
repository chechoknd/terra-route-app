package http

import "net/http"

func RegisterRoutes(mux *http.ServeMux, handler *Handler) {
	mux.HandleFunc("POST /api/v1/auth/login", handler.Login)
	mux.Handle("GET /api/v1/auth/me", AuthMiddleware(handler.tokens)(http.HandlerFunc(handler.Me)))
}
