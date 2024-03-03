package sessionid

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/crypto"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
)

type ID string

type Key struct{}

const cookieName = "session_id"

func newID() (ID, error) {
	id, err := crypto.RandString(16)
	if err != nil {
		return "", err
	}
	return ID(id), nil
}

func LoadMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		logger.Debug().Msg("Checking session status.")

		cookie, err := r.Cookie(cookieName)
		var sessionID ID
		if err != nil {
			logger.Info().Msg("Session not found. Attempting to refresh/create a new session.")

			newCookie, newID, err := getRefreshedCookie()
			if err != nil {
				logger.Error().Err(err).Msg("Failed to refresh session")
				http.Error(w, "error: Failed to refresh session", http.StatusInternalServerError)
				return
			}

			*logger = logger.With().Str("sessionID", string(newID)).Logger()
			logger.Debug().Msg("New session created and cookie set.")
			http.SetCookie(w, newCookie)
			sessionID = newID
		} else {
			sessionID = ID(cookie.Value)
			*logger = logger.With().Str("sessionID", string(sessionID)).Logger()
			logger.Debug().Msg("Existing session loaded from cookie.")
		}

		ctx := context.WithValue(r.Context(), Key{}, sessionID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRefreshedCookie() (*http.Cookie, ID, error) {
	newID, err := newID()
	if err != nil {
		return nil, "", errors.New("error: Failed to create session ID")
	}
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    string(newID),
		Path:     "/",
		Expires:  time.Now().Add(60 * time.Minute),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	return cookie, newID, nil
}

func RefreshSession(w http.ResponseWriter, r *http.Request) (ID, error) {
	newCookie, newID, err := getRefreshedCookie()
	if err != nil {
		return "", errors.New("error: Failed to refresh session")
	}
	http.SetCookie(w, newCookie)
	return newID, nil
}
