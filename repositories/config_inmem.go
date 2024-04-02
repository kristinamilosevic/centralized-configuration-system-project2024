package repositories

import (
	"errors"
	"projekat/model"
)

type ConfigInMemRepository struct {
	configs map[string]model.Config
}

func NewConfigInMemRepository() model.ConfigRepository {
	return &ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}

func (repo *ConfigInMemRepository) Create(config model.Config) error {
	if _, exists := repo.configs[config.Name]; exists {
		return errors.New("config with this name already exists")
	}

	repo.configs[config.Name] = config
	return nil
}

func (repo *ConfigInMemRepository) ReadByName(name string) (model.Config, error) {
	config, exists := repo.configs[name]
	if !exists {
		return model.Config{}, errors.New("config not found")
	}

	return config, nil
}

func (repo *ConfigInMemRepository) Update(config model.Config) error {
	if _, exists := repo.configs[config.Name]; !exists {
		return errors.New("config not found")
	}

	repo.configs[config.Name] = config
	return nil
}

func (repo *ConfigInMemRepository) DeleteByName(name string) error {
	if _, exists := repo.configs[name]; !exists {
		return errors.New("config not found")
	}

	delete(repo.configs, name)
	return nil
}
