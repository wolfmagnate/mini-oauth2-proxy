package main

import (
	"errors"
	"strings"

	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/headerInjection"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/oidc"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/upstream"
)

type ConfigSchema struct {
	OIDC            oidc.ConfigSchema            `json:"oidc"`
	Upstream        upstream.ConfigSchema        `json:"upstream"`
	HeaderInjection headerInjection.ConfigSchema `json:"headerInjection"`
	ProxyURL        proxyURL.ConfigSchema        `json:"proxyURL"`
	Log             log.ConfigSchema             `json:"log"`
	Port            int                          `json:"port" env:"OAUTH2PROXY_PORT"`
}

func (s *ConfigSchema) Validate() error {
	errMessages := make([]string, 0)

	if err := s.OIDC.Validate(); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if err := s.Upstream.Validate(); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if err := s.HeaderInjection.Validate(); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if err := s.ProxyURL.Validate(); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if err := s.Log.Validate(); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if !isValidPort(s.Port) {
		errMessages = append(errMessages, "error: port number is invalid")
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func isValidPort(p int) bool {
	return 0 <= p && p <= 65535
}

func (s *ConfigSchema) CreateConfig() any {
	return Config{
		OIDC:            s.OIDC.CreateConfig(),
		Upstream:        s.Upstream.CreateConfig(),
		HeaderInjection: s.HeaderInjection.CreateConfig(),
		ProxyURL:        s.ProxyURL.CreateConfig(),
		Log:             s.Log.CreateConfig(),
		Port:            s.Port,
	}
}
