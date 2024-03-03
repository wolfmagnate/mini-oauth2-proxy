package health

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
)

const path string = "/oauth2/health"

func AddEndpoint(r *chi.Mux) {
	r.HandleFunc(path, handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
	logger.Debug().Msg("Returned health response OK")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
	}
}
