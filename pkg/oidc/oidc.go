package oidc

import (
	"context"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/crypto"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/redirect"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/session"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/sessionid"
	"golang.org/x/oauth2"
)

const Path string = "/oauth2"

func NewRouter(config Config) *chi.Mux {
	r := chi.NewRouter()
	for _, provider := range config.providers {
		r.Handle(provider.StartPath, createOIDCStartHandler(provider))
		r.Handle(getCallbackPath(provider), createOIDCCallbackHandler(provider))
	}
	return r
}
func createOIDCStartHandler(provider Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		id := r.Context().Value(sessionid.Key{}).(sessionid.ID)
		logger.Debug().Msg("Starting OIDC authentication process")

		redirectValue := r.Context().Value(redirect.Key{})

		// 新規ログインなので、古い情報は削除してよい
		logoutCompletely(id)

		var redirectURL string
		if redirectValue != nil {
			redirectURL = redirectValue.(string)
		} else {
			// 通常、子のハンドラへのアクセスはUpstreamへのアクセスがリダイレクトされる形で行われる
			// しかし、理論的には直接OIDC認証を始めるためのエンドポイントをたたくことも出来る
			logger.Error().Msg("Authentication request without upstream redirect URL - denied")
			http.Error(w, "Authentication request without upstream redirect URL is not allowed", http.StatusBadRequest)
			return
		}

		state, err := createState()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create state for OIDC authentication")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		nonce, err := createNonce()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to create nonce for OIDC authentication")
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		session.SetState(id, state)
		session.SetNonce(id, nonce)
		session.SetRedirectURL(id, redirectURL)

		authEndpointURL := provider.OAuth2Config.AuthCodeURL(
			state,
			oidc.Nonce(nonce),
			// OAuth2.0ではredirect_uriの指定はOPTIONALだが、
			// oauth2proxyは複数のredirect_uriが使われるIdPと通信することがあるため必須である
			oauth2.SetAuthURLParam("redirect_uri", provider.OAuth2Config.RedirectURL),
		)

		logger.Info().
			Str("state", state).
			Str("nonce", nonce).
			Msg("Redirect to OIDC provider's authorization endpoint")

		http.Redirect(w, r, authEndpointURL, http.StatusFound)
	}
}
func createNonce() (string, error) {
	return crypto.RandString(16)
}

func createState() (string, error) {
	return crypto.RandString(16)
}

func getCallbackPath(provider Provider) string {
	return strings.TrimPrefix(proxyURL.GetPathFromURL(provider.OAuth2Config.RedirectURL), Path)
}

func createOIDCCallbackHandler(provider Provider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		id := r.Context().Value(sessionid.Key{}).(sessionid.ID)

		logger.Debug().Msg("Starting OIDC callback process")

		// OIDCの仕様により、StateをNonceは使ったらすぐに破棄する
		// RedirectURLも保持しておく理由がないため過ぎに破棄する
		defer func() {
			session.DeleteState(id)
			session.DeleteNonce(id)
			session.DeleteRedirectURL(id)
		}()

		// 再びログインを行おうとしているので、古いログイン情報は削除する
		// StateとNonceとRedirectURLは上のdefer節で消してくれるためlogoutでは消さなくてよい
		logout(id)

		state, err := session.GetState(id)
		if err != nil {
			logger.Error().Err(err).Msg("State not found during OIDC callback")
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("state") != state {
			logger.Error().Msg("State did not match during OIDC callback")
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := provider.OAuth2Config.Exchange(context.Background(), r.URL.Query().Get("code"))
		if err != nil {
			logger.Error().Err(err).Msg("Failed to exchange token during OIDC callback")
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			logger.Error().Msg("No id_token field in oauth2 token during OIDC callback")
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}

		idToken, err := provider.Verifier.Verify(context.Background(), rawIDToken)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to verify ID Token during OIDC callback")
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		nonce, err := session.GetNonce(id)
		if err != nil {
			logger.Error().Err(err).Msg("Nonce not found during OIDC callback")
			http.Error(w, "nonce not found", http.StatusBadRequest)
			return
		}
		if idToken.Nonce != nonce {
			logger.Error().Msg("Nonce did not match during OIDC callback")
			http.Error(w, "nonce did not match", http.StatusBadRequest)
			return
		}

		userInfo, err := provider.OIDCProvider.UserInfo(context.Background(), oauth2.StaticTokenSource(oauth2Token))
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get userInfo during OIDC callback")
			http.Error(w, "Failed to get userInfo: "+err.Error(), http.StatusInternalServerError)
			return
		}

		redirectURL, err := session.GetRedirectURL(id)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to get redirect URL during OIDC callback")
			http.Error(w, "Failed to get redirect URL: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// セッションハイジャックを防ぐため、ログインに成功したらセッションIDを再発行する
		newID, err := sessionid.RefreshSession(w, r)
		session.RefreshSession(id, newID)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to refresh session during OIDC callback")
			http.Error(w, "Failed to refresh session: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// この処理を最後に置いているのは、不正なログインを防ぐため
		// 正しくログインを検証できたときのみセッションに情報を保持する
		session.SetIDToken(newID, idToken)
		session.SetUserInfo(newID, userInfo)

		logger.Info().Msg("OIDC callback process completed successfully")
		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

func logoutCompletely(id sessionid.ID) {
	logout(id)
	session.DeleteNonce(id)
	session.DeleteRedirectURL(id)
	session.DeleteState(id)
}

func logout(id sessionid.ID) {
	session.DeleteIDToken(id)
	session.DeleteUserInfo(id)
}
