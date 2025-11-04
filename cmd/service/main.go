package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/shivaacodes/platform-service/internal/cache"
	"github.com/shivaacodes/platform-service/internal/config"
	"github.com/shivaacodes/platform-service/internal/metrics"
)

// DataResponse represents the payload returned by /api/v1/data.
type DataResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Now     string `json:"now"`
}

// dataHandler returns JSON and uses Redis as a cache-aside store.
func dataHandler(cacheClient cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		val, err := cacheClient.Get("data_key")
		if err != nil {
			resp := DataResponse{
				ID:      "sample-1",
				Message: "hello from platform-service",
				Now:     time.Now().UTC().Format(time.RFC3339),
			}
			b, _ := json.Marshal(resp)
			cacheClient.Set("data_key", string(b), time.Minute)
			val = string(b)
			w.Header().Set("X-Cache", "MISS")
		} else {
			w.Header().Set("X-Cache", "HIT")
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))

		// record metrics
		duration := time.Since(start).Seconds()
		metrics.RequestCount.WithLabelValues("/api/v1/data", r.Method, "200").Inc()
		metrics.RequestDuration.WithLabelValues("/api/v1/data", r.Method).Observe(duration)

		log.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("cache", w.Header().Get("X-Cache")).
			Dur("latency", time.Since(start)).
			Msg("request complete")
	}
}

func main() {
	// Load configuration
	config.Load()
	cfg := config.C

	metrics.Init()

	// Initialize Redis
	cacheClient := cache.NewClient(cfg.RedisAddr, cfg.RedisPassword)
	defer cacheClient.Close()

	// Setup HTTP routes
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Str("method", r.Method).Str("path", r.URL.Path).Msg("health check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/readyz", ReadyzHandler(cacheClient))
	mux.HandleFunc("/api/v1/data", dataHandler(cacheClient))
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in background
	go func() {
		log.Info().Str("port", cfg.Port).Msg("server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info().Msg("shutdown signal received")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	} else {
		log.Info().Msg("server exited cleanly")
	}
}
