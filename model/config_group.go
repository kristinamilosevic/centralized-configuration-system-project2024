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
	// dodati metode (crud)
}
