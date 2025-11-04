package logging

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init configures global logging behavior using zerolog.
// - Uses human-readable console logs in development.
// - Uses structured JSON logs in production.
// - Adds timestamps for consistent event tracking.

func Init(env string) {
	// Set timestamp format to Unix time for compact, standard logs
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Configure output format based on environment

	if env == "development" {
		// Pretty console output for local readability
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	} else {
		// JSON logs with timestamps for production observability and log aggregation
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
