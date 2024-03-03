package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/config"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/headerInjection"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/health"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/login"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/oidc"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/ready"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/redirect"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/requestid"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/session"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/sessionid"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/upstream"
)

func main() {
	StartOAuth2Proxy()
}

func StartOAuth2Proxy() {
	c := config.LoadConfig(&ConfigSchema{}).(Config)
	headerInjectMiddleware := headerInjection.CreateMiddleware(c.HeaderInjection)
	proxyURL.Init(c.ProxyURL)
	session.Init()
	log.Init(c.Log)

	r := chi.NewRouter()
	r.Use(log.CreateLoggerMiddleware)
	r.Use(requestid.AddIDMiddleware)
	r.Use(sessionid.LoadMiddleware)
	r.Use(login.GetLoginStatusMiddleware)
	r.Use(redirect.GetMiddleware)
	health.AddEndpoint(r)
	ready.AddEndpoint(r)
	oidcRouter := oidc.NewRouter(c.OIDC)
	r.Mount(oidc.Path, oidcRouter)

	upstreamRouter := upstream.NewRouter(c.Upstream)
	loginHandler := oidc.NewLoginHandler(c.OIDC)

	r.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
		isLogin := r.Context().Value(login.Key{}).(bool)
		if isLogin {
			headerInjectMiddleware(upstreamRouter).ServeHTTP(w, r)
		} else {
			redirect.FindMiddleware(loginHandler).ServeHTTP(w, r)
		}
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r); err != nil {
		panic(err.Error())
	}
}
