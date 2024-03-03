package log

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Key struct{}

var config Config

func Init(c Config) {
	config = c
}

func CreateLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zerolog.TimeFieldFormat = time.RFC3339
		logger := log.Output(os.Stdout).Level(config.Level)
		ctx := context.WithValue(r.Context(), Key{}, &logger)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
