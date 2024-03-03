package upstream

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
)

func NewRouter(config Config) *chi.Mux {
	r := chi.NewRouter()
	for _, server := range config.Servers {
		proxy := setupReverseProxy(server)
		r.Route(server.MatchPrefix, func(r chi.Router) {
			r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
				logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
				logger.Info().Msg(fmt.Sprintf("Proxied request to upstream: %s.", server.ID))
				proxy.ServeHTTP(w, r)
			})
		})
	}
	return r
}

func setupReverseProxy(server Server) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(server.URL)
	proxy.Director = modifyRequest(server)
	setTimeout(server, proxy)
	proxy.ModifyResponse = modifyResponse(server)
	return proxy
}

func setTimeout(server Server, proxy *httputil.ReverseProxy) {
	if server.Timeout != nil {
		proxy.Transport = &http.Transport{
			ResponseHeaderTimeout: time.Duration(*server.Timeout),
		}
	}
}
