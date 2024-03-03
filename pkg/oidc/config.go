package oidc

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type Config struct {
	providers     []Provider
	skipLoginPage bool
}

type Provider struct {
	ID           string
	StartPath    string
	OIDCProvider *oidc.Provider
	Verifier     *oidc.IDTokenVerifier
	OAuth2Config *oauth2.Config
}
