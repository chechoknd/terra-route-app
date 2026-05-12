package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	authusecases "github.com/terraroute/terra-route/backend/internal/auth/application/usecases"
	authinfra "github.com/terraroute/terra-route/backend/internal/auth/infrastructure"
	authhttp "github.com/terraroute/terra-route/backend/internal/auth/interfaces/http"
	userinfra "github.com/terraroute/terra-route/backend/internal/users/infrastructure"
)

type Config struct {
	Addr         string
	DB           *pgxpool.Pool
	Logger       *slog.Logger
	JWTSecret    string
	JWTExpiresIn string
}

type Server struct {
	httpServer *http.Server
	db         *pgxpool.Pool
	logger     *slog.Logger
}

func New(cfg Config) (*Server, error) {
	mux := http.NewServeMux()

	srv := &Server{
		db:     cfg.DB,
		logger: cfg.Logger,
	}

	userRepo := userinfra.NewPostgresRepository(cfg.DB)
	passwordHasher := authinfra.NewBcryptPasswordHasher(0)
	tokenService, err := authinfra.NewJWTService(cfg.JWTSecret, cfg.JWTExpiresIn)
	if err != nil {
		return nil, fmt.Errorf("create jwt service: %w", err)
	}
	loginUseCase := authusecases.NewLoginUseCase(userRepo, passwordHasher, tokenService)
	authHandler := authhttp.NewHandler(loginUseCase, tokenService, userRepo)

	mux.HandleFunc("GET /healthz", srv.handleHealth)
	mux.HandleFunc("GET /readyz", srv.handleReady)
	mux.HandleFunc("GET /api/v1/healthz", srv.handleHealth)
	authhttp.RegisterRoutes(mux, authHandler)

	srv.httpServer = &http.Server{
		Addr:              cfg.Addr,
		Handler:           requestLogger(cfg.Logger, mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	return srv, nil
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "terraroute-api",
	})
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := s.db.Ping(ctx); err != nil {
		s.logger.Warn("readiness check failed", slog.String("error", err.Error()))
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{
			"status": "unavailable",
		})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func requestLogger(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		logger.Info("http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", time.Since(start)),
		)
	})
}
