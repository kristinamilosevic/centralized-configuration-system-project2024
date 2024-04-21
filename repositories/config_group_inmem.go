package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type ConfigGroupInMemRepository struct {
	configGroups map[string]model.ConfigGroup
}

func NewConfigGroupInMemRepository() model.ConfigGroupRepository {
	return &ConfigGroupInMemRepository{
		configGroups: make(map[string]model.ConfigGroup),
	}
}

func (repo *ConfigGroupInMemRepository) Create(configGroup model.ConfigGroup) error {
	if _, exists := repo.configGroups[configGroup.Name]; exists {
		return errors.New("config group with this name already exists")
	}

	repo.configGroups[configGroup.Name] = configGroup
	return nil
}

func (repo *ConfigGroupInMemRepository) Read(name string, version int) (model.ConfigGroup, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	configGroup, exists := repo.configGroups[key]
	if !exists {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return configGroup, nil
}

func (repo *ConfigGroupInMemRepository) Update(newConfigGroup model.ConfigGroup) error {
	key := configGroupKey(newConfigGroup.Name, newConfigGroup.Version)
	if _, exists := repo.configGroups[key]; !exists {
		return errors.New("config group not found")
	}
	repo.configGroups[key] = newConfigGroup
	return nil
}

func (repo *ConfigGroupInMemRepository) Delete(name string, version int) error {
	key := fmt.Sprintf("%s/%d", name, version)
	if _, exists := repo.configGroups[key]; !exists {
		return errors.New("config group not found")
	}
	delete(repo.configGroups, key)
	return nil
}

// GetAll vraća sve konfiguracije
func (repo *ConfigGroupInMemRepository) GetAll() ([]model.ConfigGroup, error) {
	configGroups := make([]model.ConfigGroup, 0, len(repo.configGroups))
	for _, configGroup := range repo.configGroups {
		configGroups = append(configGroups, configGroup)
	}
	return configGroups, nil
}

func (c ConfigGroupInMemRepository) Add(configGroup model.ConfigGroup) {
	key := fmt.Sprintf("%s/%d", configGroup.Name, configGroup.Version)
	c.configGroups[key] = configGroup
}

// configKey kreira ključ za konfiguraciju na osnovu imena i verzije
func configGroupKey(name string, version int) string {
	return fmt.Sprintf("%s/%d", name, version)
}

func (c ConfigGroupInMemRepository) Get(name string, version int) (model.ConfigGroup, error) {
	key := configGroupKey(name, version)
	configGroup, ok := c.configGroups[key]
	if !ok {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	return configGroup, nil
}
