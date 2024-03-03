package proxyURL

type ConfigSchema struct {
	Host string `json:"host"`
}

func (s *ConfigSchema) Validate() error {
	return nil
}

func (s *ConfigSchema) CreateConfig() Config {
	return Config{
		Host: s.Host,
	}
}
