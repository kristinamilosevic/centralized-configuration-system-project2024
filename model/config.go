package model

type Config struct {
	Name       string            `json:"name"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"parameters"`
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
	Add(Config Config)
	Get(name string, version int) (Config, error)
	GetAll() ([]Config, error)
}
