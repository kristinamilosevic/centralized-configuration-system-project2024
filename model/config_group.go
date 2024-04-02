package model

type ConfigGroup struct {
	Name          string
	Version       int
	Configuration []Config
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
}
