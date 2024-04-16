package model

type ConfigGroup struct {
	Name          string   `json:"name"`
	Version       int      `json:"version"`
	Configuration []Config `json:"configuration"`
}

func NewConfigGroup(name string, version int, configuration []Config) ConfigGroup {
	return ConfigGroup{
		Name:          name,
		Version:       version,
		Configuration: configuration,
	}
}

type ConfigGroupRepository interface {
	Create(configGroup ConfigGroup) error
	ReadByName(name string) (ConfigGroup, error)
	Update(configGroup ConfigGroup) error
	DeleteByName(name string) error
	GetAll() ([]ConfigGroup, error)
	Add(ConfigGroup ConfigGroup)
	Get(name string, version int) (ConfigGroup, error)
}
