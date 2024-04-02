package repositories

import (
	"errors"
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

func (repo *ConfigGroupInMemRepository) ReadByName(name string) (model.ConfigGroup, error) {
	configGroup, exists := repo.configGroups[name]
	if !exists {
		return model.ConfigGroup{}, errors.New("config group not found")
	}

	return configGroup, nil
}

func (repo *ConfigGroupInMemRepository) Update(configGroup model.ConfigGroup) error {
	if _, exists := repo.configGroups[configGroup.Name]; !exists {
		return errors.New("config group not found")
	}

	repo.configGroups[configGroup.Name] = configGroup
	return nil
}

func (repo *ConfigGroupInMemRepository) DeleteByName(name string) error {
	if _, exists := repo.configGroups[name]; !exists {
		return errors.New("config group not found")
	}

	delete(repo.configGroups, name)
	return nil
}
