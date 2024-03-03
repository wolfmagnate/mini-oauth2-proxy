package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type ConfigSchema struct {
	Providers     []ProviderSchema `json:"providers"`
	SkipLoginPage bool             `json:"skipLoginPage"`
}

type ProviderSchema struct {
	ID           string `json:"id"`
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	RedirectURL  string `json:"redirectURL"`
	StartPath    string `json:"startPath"`

	// スコープの内、"oidc"を除いたもの。oidcは自動追加するため不要
	Scopes []string `json:"scopes"`

	// IssuerからOIDC Discoveryを使うため、その他の情報は不要
	Issuer string `json:"issuer"`
}

func (s *ConfigSchema) Validate() error {
	errMessages := make([]string, 0)

	if len(s.Providers) < 1 {
		errMessages = append(errMessages, "error: at least one provider is required")
	}

	providerIDs := make(map[string]bool)
	for _, p := range s.Providers {
		if _, exists := providerIDs[p.ID]; exists {
			errMessages = append(errMessages, "error: duplicate upstream ID")
		}
		providerIDs[p.ID] = true
	}

	if err := validateURLs(s.Providers); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if err := validateStartPath(s.Providers); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if s.SkipLoginPage && len(s.Providers) > 1 {
		errMessages = append(errMessages, "error: cannot skip login page because there are more than one provider")
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func validateURLs(providers []ProviderSchema) error {
	errMessages := make([]string, 0)

	for _, p := range providers {
		if !strings.HasPrefix(p.Issuer, "https://") || !isValidURL(p.Issuer) {
			errMessages = append(errMessages, fmt.Sprintf("error: provider issuer is not a valid https URL: %s", p.Issuer))
		}
		if !strings.HasPrefix(p.RedirectURL, "https://") || !isValidURL(p.RedirectURL) {
			errMessages = append(errMessages, fmt.Sprintf("error: provider redirectURL is not a valid https URL: %s", p.RedirectURL))
		}
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func validateStartPath(providers []ProviderSchema) error {
	errMessages := make([]string, 0)

	for _, p := range providers {
		if !isValidPath(p.StartPath) {
			errMessages = append(errMessages, fmt.Sprintf("error: provider startPath is not a valid path: %s", p.StartPath))
		}
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func isValidURL(toTest string) bool {
	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func isValidPath(path string) bool {
	regex := regexp.MustCompile(`^/([A-Za-z0-9\-_]+(/|$))*$`)
	return regex.MatchString(path)
}

func (s *ConfigSchema) CreateConfig() Config {
	providers := make([]Provider, 0)
	for _, p := range s.Providers {
		ctx := context.Background()
		provider, err := oidc.NewProvider(ctx, p.Issuer)
		if err != nil {
			panic(err)
		}

		oidcConfig := &oidc.Config{
			ClientID: p.ClientID,
		}
		verifier := provider.Verifier(oidcConfig)

		scopes := p.Scopes
		scopes = append(scopes, oidc.ScopeOpenID)
		config := &oauth2.Config{
			ClientID:     p.ClientID,
			ClientSecret: p.ClientSecret,
			Endpoint:     provider.Endpoint(),
			RedirectURL:  p.RedirectURL,
			Scopes:       scopes,
		}

		providers = append(providers, Provider{
			ID:           p.ID,
			StartPath:    p.StartPath,
			OIDCProvider: provider,
			Verifier:     verifier,
			OAuth2Config: config,
		})
	}
	return Config{
		providers:     providers,
		skipLoginPage: s.SkipLoginPage,
	}
}
