package login

import (
	"context"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/session"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/sessionid"
)

type Key struct{}

func GetLoginStatusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(sessionid.Key{}).(sessionid.ID)
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)

		logger.Debug().Msg("Checking login status.")

		_, err := session.GetIDToken(id)
		var ctx context.Context
		if err == nil {
			logger.Debug().Msg("User is logged in.")
			*logger = logger.With().Bool("login", true).Logger()
			ctx = context.WithValue(r.Context(), Key{}, true)
		} else if err.Error() == "error: IDToken not found" {
			logger.Debug().Msg("User is not logged in.")
			*logger = logger.With().Bool("login", false).Logger()
			ctx = context.WithValue(r.Context(), Key{}, false)
		} else {
			logger.Error().Err(err).Msg("Unexpected error while fetching login status.")
			http.Error(w, "error: unexpected error while fetching login status", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
