package oidc

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/redirect"
)

//go:embed login_page.html
var loginPageHTML embed.FS

type templateData struct {
	Providers []struct {
		ID          string
		LoginURL    string
		RedirectURL string
	}
}

func NewLoginHandler(config Config) http.Handler {
	return noCacheMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		upstreamRedirectURL := r.Context().Value(redirect.Key{}).(string)

		logger.Debug().Msg("Handling new login request")

		if config.skipLoginPage {
			logger.Info().Msg("Skipping login page and redirecting to IdP")
			baseURL := proxyURL.GetURLFromPath(Path + config.providers[0].StartPath)
			query := url.Values{}
			query.Set("redirect", upstreamRedirectURL)
			baseURL.RawQuery = query.Encode()
			http.Redirect(w, r, baseURL.String(), http.StatusFound)
			return
		}

		logger.Debug().Msg("Rendering login page")

		tmpl, err := template.ParseFS(loginPageHTML, "login_page.html")
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse login page template")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		data := getTemplateData(config.providers, upstreamRedirectURL)

		err = tmpl.Execute(w, data)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to execute login page template")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		logger.Info().Msg("Login page rendered successfully")
	}))
}

func getTemplateData(providers []Provider, upstreamRedirectURL string) templateData {
	var providersData []struct {
		ID          string
		LoginURL    string
		RedirectURL string
	}
	for _, provider := range providers {
		baseURL := proxyURL.GetURLFromPath(Path + provider.StartPath)
		loginURL := baseURL.String()
		providersData = append(providersData, struct {
			ID          string
			LoginURL    string
			RedirectURL string
		}{
			ID:          provider.ID,
			LoginURL:    loginURL,
			RedirectURL: upstreamRedirectURL,
		})
	}
	return templateData{
		Providers: providersData,
	}
}
