package model

type Config2 struct {
	Name    string            `json:"name"`
	Version int               `json:"version"`
	Labels  map[string]string `json:"labels"`
}

func NewConfig_2(name string, version int, labels map[string]string) Config2 {
	return Config2{
		Name:    name,
		Version: version,
		Labels:  labels,
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
	Add(Config Config2)
	Get(name string, version int) (Config2, error)
	GetAll() ([]Config2, error)
}
