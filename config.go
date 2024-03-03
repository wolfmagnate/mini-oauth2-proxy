package main

import (
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/headerInjection"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/oidc"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/upstream"
)

type Config struct {
	OIDC            oidc.Config
	Upstream        upstream.Config
	HeaderInjection headerInjection.Config
	ProxyURL        proxyURL.Config
	Log             log.Config
	Port            int
}
