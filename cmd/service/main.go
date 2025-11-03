package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shivaacodes/platform-service/internal/cache"
	"github.com/shivaacodes/platform-service/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type DataResponse struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	Now     string `json:"now"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// init metrics and redis
	metrics.Init()
	cacheClient := cache.NewClient("redis:6379")
	defer cacheClient.Close()

	mux := http.NewServeMux()

	// health endpoint
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("health check: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// data endpoint with cache + metrics
	mux.HandleFunc("/api/v1/data", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("request: %s %s", r.Method, r.URL.Path)

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
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))

		duration := time.Since(start).Seconds()
		metrics.RequestCount.WithLabelValues("/api/v1/data", "GET", "200").Inc()
		metrics.RequestDuration.WithLabelValues("/api/v1/data", "GET").Observe(duration)
	})

	// metrics endpoint
	mux.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// graceful shutdown
	go func() {
		log.Printf("Server running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down...")
	server.Shutdown(context.Background())
}

