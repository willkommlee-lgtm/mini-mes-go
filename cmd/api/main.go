// Command api starts the mini-mes-go HTTP server.
//
// Day 1: only a health check. The real endpoints (work orders, lots, scans)
// are added on Day 2 — see README.md for the roadmap.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/willkommlee-lgtm/mini-mes-go/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) // structured JSON logs (useful again on Day 5)

	addr := getenv("HTTP_ADDR", ":8080")

	// Connect to PostgreSQL
	dbURL := getenv("DATABASE_URL", "postgres://mes:mes@localhost:5432/mes")
	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("invalid database config", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		logger.Error("database not reachable", "err", err)
		os.Exit(1)
	}
	logger.Info("database connected")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	scans := &scanHandler{materials: store.NewPostgres(db), logger: logger}
	mux.HandleFunc("POST /scans", scans.handleScan)

	// TODO(Day 2): connect to PostgreSQL (DATABASE_URL) and register:
	//   POST /work-orders, GET /work-orders/{id}
	//   POST /lots
	//   POST /scans                -> validate a material barcode (ALLOW / BLOCK)
	//   GET  /lots/{id}/genealogy  -> which materials went into a lot

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in a goroutine so main can wait for a shutdown signal.
	go func() {
		logger.Info("server starting", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown on Ctrl+C or SIGTERM (Kubernetes sends SIGTERM on Day 6).
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown failed", "err", err)
	}
	logger.Info("server stopped")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
