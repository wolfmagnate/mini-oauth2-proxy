package proxyURL

import "net/url"

var config Config

func Init(c Config) {
	config = c
}

func GetURLFromPath(path string) *url.URL {
	return &url.URL{
		Scheme: "https",
		Host:   config.Host,
		Path:   path,
	}
}

func GetPathFromURL(URL string) string {
	parsedURL, _ := url.Parse(URL)
	return parsedURL.Path
}
