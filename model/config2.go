package model

type Config2 struct {
	Name       string            `json:"name"`
	Version    int               `json:"version"`
	Parameters map[string]string `json:"parameters"`
	Labels     map[string]string `json:"labels"`
}

func NewConfig_2(name string, version int, parameters map[string]string, labels map[string]string) Config2 {
	return Config2{
		Name:       name,
		Version:    version,
		Parameters: parameters,
		Labels:     labels,
	}
}

func NewConfig2_2(name string, version int) Config2 {
	return Config2{
		Name:    name,
		Version: version,
	}
}

type Config2Repository interface {
	Create(config Config2) error
	Read(name string, version int) (Config2, error)
	Update(config Config2) error
	Delete(name string, version int) error
	Get(name string, version int) (Config2, error)
	GetAll() ([]Config2, error)
}
