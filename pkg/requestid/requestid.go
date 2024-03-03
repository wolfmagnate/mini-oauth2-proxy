package requestid

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/crypto"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
)

type ID string

type Key struct{}

func newID() (ID, error) {
	id, err := crypto.RandString(16)
	if err != nil {
		return "", err
	}
	return ID(id), nil
}

func AddIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		logger.Debug().Msg("Adding request ID")

		id, err := newID()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create request ID")
			http.Error(w, "error: Failed to create request ID", http.StatusInternalServerError)
			return
		}

		*logger = logger.With().Str("requestID", string(id)).Logger()
		ctx := context.WithValue(r.Context(), Key{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
