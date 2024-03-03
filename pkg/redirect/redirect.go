package redirect

import (
	"context"
	"net/http"
	"net/url"

	"github.com/rs/zerolog"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/log"
	"github.com/wolfmagnate/mini-oauth2-proxy/pkg/proxyURL"
)

type Key struct{}

// 指定されたリクエストから考えられる認証成功後のリダイレクトURLを探索する
// 明示的に指定されていないものでも、リクエスト情報をもとに決定するため、Upstreamへのリクエストにしか使用できない
func FindMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		logger.Debug().Msg("Finding application redirect URL.")
		url := getRedirectURL(r)
		*logger = logger.With().Str("found redirect URL", url).Logger()
		ctx := context.WithValue(r.Context(), Key{}, url)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// 指定されたリクエストに明示的に指定されている認証成功後のリダイレクトURLを取得する
// OIDC用エンドポイントへの直接のアクセスなどプロキシ先ではないoauth2proxy自体へのアクセスの場合、
// 暗黙的にリダイレクトURLを決定することが出来ないため、この処理が必要
func GetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := r.Context().Value(log.Key{}).(*zerolog.Logger)
		logger.Debug().Msg("Checking specified application redirect URL.")

		url := getExplicitRedirectURL(r)
		if url != "" {
			logger.Debug().Msg("application redirect URL was found")
			*logger = logger.With().Str("specified redirect URL", url).Logger()
		} else {
			logger.Debug().Msg("application redirect URL was not found")
		}
		ctx := context.WithValue(r.Context(), Key{}, url)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getExplicitRedirectURL(r *http.Request) string {
	if canRedirectByQueryParam(r.URL.Query()) {
		return getURLFromQueryParam(r.URL.Query())
	}
	if canRedirectByHeader(r.Header) {
		return getURLFromHeader(r.Header)
	}
	return ""
}

func getRedirectURL(r *http.Request) string {
	if explicitURL := getExplicitRedirectURL(r); explicitURL != "" {
		return explicitURL
	}
	return proxyURL.GetURLFromPath(r.URL.Path).String()
}

func canRedirectByQueryParam(params url.Values) bool {
	redirectURL := params.Get("redirect")
	return isValidURL(redirectURL)
}

func canRedirectByHeader(header http.Header) bool {
	redirectURL := header.Get("X-Auth-Request-Redirect")
	return isValidURL(redirectURL)
}

func getURLFromQueryParam(params url.Values) string {
	return params.Get("redirect")
}

func getURLFromHeader(header http.Header) string {
	return header.Get("X-Auth-Request-Redirect")
}

func isValidURL(toTest string) bool {
	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != ""
}
