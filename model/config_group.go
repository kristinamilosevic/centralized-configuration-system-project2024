package model

// swagger:model ConfigGroup
type ConfigGroup struct {
	// Name of the post
	// in: string
	Name string `json:"name"`
	// Version of the post
	// in: int
	Version int `json:"version"`
	// Configuration of the post
	// in: []Config2
	Configuration []Config2 `json:"configuration"`
}

func NewConfigGroup(name string, version int, configuration []Config2) ConfigGroup {
	return ConfigGroup{
		Name:          name,
		Version:       version,
		Configuration: configuration,
	}
}

func NewConfigGroup2(name string, version int) ConfigGroup {
	return ConfigGroup{
		Name:    name,
		Version: version,
	}
}

type ConfigGroupRepository interface {
	CreateConfigGroup(configGroup *ConfigGroup, idempotencyKey, bodyHash string) error
	CheckIfExists(idempotencyKey, bodyHash string) (bool, error)
	Read(name string, version int) (ConfigGroup, error)
	Update(configGroup ConfigGroup) error
	Delete(name string, version int) error
	GetAll() ([]ConfigGroup, error)
	Get(name string, version int) (ConfigGroup, error)
	RemoveConfig(groupName string, groupVersion int, configName string, configVersion int) error
	AddConfig(groupName string, groupVersion int, config Config2) error
	GetFilteredConfigs(name string, version int, filter map[string]string) ([]Config2, error)
	RemoveByLabels(groupName string, groupVersion int, filter map[string]string) error
}
