package log

import (
	"errors"

	"github.com/rs/zerolog"
)

type ConfigSchema struct {
	Level string `json:"level" env:"OAUTH2PROXY_LOG_LEVEL" envDefault:"Info"`
}

func (s *ConfigSchema) Validate() error {
	switch s.Level {
	case "Debug", "Info", "Warn", "Error":
		return nil
	default:
		return errors.New("error: invalid log level")
	}
}

func (s *ConfigSchema) CreateConfig() Config {
	var level zerolog.Level
	switch s.Level {
	case "Debug":
		level = zerolog.DebugLevel
	case "Info":
		level = zerolog.InfoLevel
	case "Warn":
		level = zerolog.WarnLevel
	case "Error":
		level = zerolog.ErrorLevel
	}
	return Config{
		Level: level,
	}
}
