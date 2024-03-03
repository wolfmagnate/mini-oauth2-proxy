package upstream

import "net/url"

type Config struct {
	Servers []Server
}

type Server struct {
	ID          string
	URL         *url.URL
	MatchPrefix string
	Timeout     *Duration
}
