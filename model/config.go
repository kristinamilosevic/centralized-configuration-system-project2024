package model

type Config struct {
	name       string
	version    int
	parameters map[string]string
	// dodati atribute
}

func NewConfig(name string, version int, parameters map[string]string) Config {
	return Config{
		name:       name,
		version:    version,
		parameters: parameters,
	}
}

type ConfigRepository interface {
	// dodati metode (crud)
}
