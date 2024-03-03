package headerInjection

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/sessionid"
)

func CreateMiddleware(config Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
			id := r.Context().Value(sessionid.Key{}).(sessionid.ID)

			logger.Debug().Msg("Injecting header of upstream request")

			for _, injector := range config.Request {
				key := injector.GetKey()
				value, err := injector.GetValue(id)
				if err != nil {
					logger.Error().Str("headerKey", key).Err(err).Msg("Failed to set request header")
					http.Error(w, fmt.Sprintf("Error setting request header '%s': %v", key, err), http.StatusInternalServerError)
					return
				}
				r.Header.Set(key, value)
				logger.Debug().Str("headerKey", key).Str("headerValue", value).Msg("Request header set successfully")
			}

			logger.Debug().Msg("Completed setting request headers")

			for _, injector := range config.Response {
				key := injector.GetKey()
				value, err := injector.GetValue(id)
				if err != nil {
					logger.Error().Str("headerKey", key).Err(err).Msg("Failed to set response header")
					http.Error(w, fmt.Sprintf("Error setting response header '%s': %v", key, err), http.StatusInternalServerError)
					return
				}
				w.Header().Set(key, value)
				logger.Debug().Str("headerKey", key).Str("headerValue", value).Msg("Response header set successfully")
			}

			logger.Debug().Msg("Completed setting response headers")

			next.ServeHTTP(w, r)
		})
	}
}
