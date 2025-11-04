package main

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/shivaacodes/platform-service/internal/cache"
)

// Checks Redis Connectivity
func ReadyzHandler(cacheClient cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 500*time.Millisecond)
		defer cancel()

		if err := cacheClient.Ping(ctx); err != nil {
			log.Error().Err(err).Msg("Redis not reachable")
			http.Error(w, "redis not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}
