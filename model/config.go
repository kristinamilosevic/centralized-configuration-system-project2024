package model

type Config struct {
	Name       string
	Version    int
	Parameters map[string]string
}

func NewConfig(name string, version int, parameters map[string]string) Config {
	return Config{
		Name:       name,
		Version:    version,
		Parameters: parameters,
	}
}

type ConfigRepository interface {
	Create(config Config) error
	ReadByName(name string) (Config, error)
	Update(config Config) error
	DeleteByName(name string) error
}
