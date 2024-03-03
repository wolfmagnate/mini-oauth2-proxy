package upstream

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ConfigSchema struct {
	Servers []ServerSchema `json:"servers"`
}

type ServerSchema struct {
	ID          string    `json:"id"`
	URL         string    `json:"url"`
	MatchPrefix string    `json:"matchPath"`
	Timeout     *Duration `json:"timeout,omitempty"`
}

// 期間をそのままJSONに記述できるようにするためには、encoding/jsonの要求するインターフェースをみたす型である必要があるため
type Duration time.Duration

func (d *Duration) UnmarshalJSON(data []byte) error {
	input := string(data)
	if unquoted, err := strconv.Unquote(input); err == nil {
		input = unquoted
	}

	du, err := time.ParseDuration(input)
	if err != nil {
		return err
	}
	*d = Duration(du)
	return nil
}

func (s *ConfigSchema) Validate() error {
	errMessages := make([]string, 0)
	if len(s.Servers) == 0 {
		errMessages = append(errMessages, "error: at least one upstream is required")
	}

	upstreamIDs := make(map[string]bool)
	for _, u := range s.Servers {
		if _, exists := upstreamIDs[u.ID]; exists {
			errMessages = append(errMessages, "error: duplicate upstream ID")
		}
		upstreamIDs[u.ID] = true
	}

	if err := validateURLs(s.Servers); err != nil {
		errMessages = append(errMessages, err.Error())
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func validateURLs(servers []ServerSchema) error {
	errMessages := make([]string, 0)

	for _, s := range servers {
		if !strings.HasPrefix(s.URL, "http://") || !isValidURL(s.URL) {
			errMessages = append(errMessages, fmt.Sprintf("error: upstream server's URL is not a valid http URL: %s", s.URL))
		}
	}

	if len(errMessages) > 0 {
		return errors.New(strings.Join(errMessages, "\n"))
	}
	return nil
}

func isValidURL(toTest string) bool {
	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (s *ConfigSchema) CreateConfig() Config {
	servers := make([]Server, 0)
	for _, server := range s.Servers {
		// Validateにてエラーチェックは終わっているため不要
		baseURL, _ := url.Parse(server.URL)

		timeout := server.Timeout
		if timeout == nil {
			t := 30 * time.Second
			timeout = (*Duration)(&t)
		}
		servers = append(servers, Server{
			ID:          server.ID,
			URL:         baseURL,
			MatchPrefix: server.MatchPrefix,
			Timeout:     timeout,
		})
	}

	return Config{
		Servers: servers,
	}
}
