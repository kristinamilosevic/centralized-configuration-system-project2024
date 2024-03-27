package model

type ConfigGroup struct {
	name          string
	version       int
	configuration []Config
	// dodati atribute
}

func NewConfigGroup(name string, version int, configuration []Config) ConfigGroup {
	return ConfigGroup{
		name:          name,
		version:       version,
		configuration: configuration,
	}
}

type ConfigGroupRepository interface {
	// dodati metode (crud)
}
