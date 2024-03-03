package headerInjection

import (
	"errors"
	"fmt"
	"strings"
)

type ConfigSchema struct {
	Request  []HeaderSchema `json:"request"`
	Response []HeaderSchema `json:"response"`
}

type HeaderSchema struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Values []string `json:"values"`
}

func (s *ConfigSchema) Validate() error {
	errMessages := make([]string, 0)

	if err := validateHeaderTypeAndValue(s.Request); err != nil {
		errMessages = append(errMessages, err.Error())
	}
	if err := validateHeaderTypeAndValue(s.Response); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if err := validateUniqueHeaderNames(s.Request); err != nil {
		errMessages = append(errMessages, err.Error())
	}
	if err := validateUniqueHeaderNames(s.Response); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func validateUniqueHeaderNames(headers []HeaderSchema) error {
	headerNames := make(map[string]bool)
	for _, header := range headers {
		if _, exists := headerNames[header.Name]; exists {
			return fmt.Errorf("error: duplicate header name found in response headers: %s", header.Name)
		}
		headerNames[header.Name] = true
	}
	return nil
}

func validateHeaderTypeAndValue(headers []HeaderSchema) error {
	validHeaderValues := map[string]bool{
		"name": true, "family_name": true, "given_name": true, "middle_name": true,
		"nickname": true, "preferred_username": true, "profile": true, "picture": true,
		"website": true, "gender": true, "birthdate": true, "zoneinfo": true,
		"locale": true, "updated_at": true, "email": true, "email_verified": true,
		"address": true, "phone_number": true, "phone_number_verified": true,
	}
	invalidHeaders := make([]string, 0)
	for _, h := range headers {
		if h.Type != "userInfo" && h.Type != "idTokenClaim" {
			invalidHeaders = append(invalidHeaders, fmt.Sprintf("error: header type is invalid: %s", h.Type))
		}
		if len(h.Values) == 0 {
			invalidHeaders = append(invalidHeaders, "error: no values to inject")
		}
		for _, v := range h.Values {
			if !validHeaderValues[v] {
				invalidHeaders = append(invalidHeaders, fmt.Sprintf("error: invalid header value: %s", v))
			}
		}
	}
	if len(invalidHeaders) > 0 {
		return errors.New(strings.Join(invalidHeaders, "\n"))
	}
	return nil
}

func (s *ConfigSchema) CreateConfig() Config {
	requestInjectors := make([]headerInjector, 0)
	responseInjectors := make([]headerInjector, 0)
	for _, header := range s.Request {
		requestInjectors = append(requestInjectors, createInjector(header))
	}
	for _, header := range s.Response {
		responseInjectors = append(responseInjectors, createInjector(header))
	}
	return Config{
		Request:  requestInjectors,
		Response: responseInjectors,
	}
}

func createInjector(s HeaderSchema) headerInjector {
	switch s.Type {
	case "userInfo":
		return createUserInfoInjector(s)
	case "idTokenClaim":
		return createIdTokenInjector(s)
	}
	panic("error: unknown HeaderSchema type")
}

func createUserInfoInjector(s HeaderSchema) *userInfoInjector {
	claims := make([]ClaimType, 0)
	for _, v := range s.Values {
		c, err := toClaimType(v)
		if err != nil {
			panic(fmt.Sprintf("error: could not convert claim %s", v))
		}
		claims = append(claims, c)
	}
	return &userInfoInjector{
		Name:   s.Name,
		Claims: claims,
	}
}

func createIdTokenInjector(s HeaderSchema) *idTokenInjector {
	claims := make([]ClaimType, 0)
	for _, v := range s.Values {
		c, err := toClaimType(v)
		if err != nil {
			panic(fmt.Sprintf("error: could not convert claim %s", v))
		}
		claims = append(claims, c)
	}
	return &idTokenInjector{
		Name:   s.Name,
		Claims: claims,
	}
}
