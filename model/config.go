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

func NewConfig2(name string, version int) Config {
	return Config{
		Name:    name,
		Version: version,
	}
}

type ConfigRepository interface {
	Create(config Config) error
	Read(name string, version int) (Config, error)
	Update(config Config) error
	Delete(name string, version int) error
	Add(Config Config)
	Get(name string, version int) (Config, error)
	GetAll() ([]Config, error)
}
