package upstream

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
)

func modifyRequest(server Server) func(r *http.Request) {
	return func(r *http.Request) {
		fixUpstreamPath(server, r)
	}
}

func fixUpstreamPath(server Server, r *http.Request) {
	originalPath := r.URL.Path
	r.URL.Scheme = server.URL.Scheme
	r.URL.Host = server.URL.Host
	r.URL.Path = server.URL.Path + strings.TrimPrefix(originalPath, server.MatchPrefix)
}

func modifyResponse(server Server) func(response *http.Response) error {
	return func(response *http.Response) error {
		if needsLocationHeader(response) {
			fixUpstreamRedirectResponsePath(server, response)
		}
		return nil
	}
}

func needsLocationHeader(response *http.Response) bool {
	return response.StatusCode == http.StatusMovedPermanently || response.StatusCode == http.StatusFound || response.StatusCode == http.StatusSeeOther || response.StatusCode == http.StatusTemporaryRedirect || response.StatusCode == http.StatusPermanentRedirect
}

func fixUpstreamRedirectResponsePath(server Server, response *http.Response) error {

	locationHeader := response.Header.Get("Location")
	if locationHeader != "" {
		parsedURL, err := url.Parse(locationHeader)
		if err != nil {
			return err
		}
		newPath := strings.TrimPrefix(parsedURL.Path, server.URL.Path)
		newPath = server.MatchPrefix + newPath
		proxiedURL := proxyURL.GetURLFromPath(newPath)
		response.Header.Set("Location", proxiedURL.String())
	}
	return nil
}
