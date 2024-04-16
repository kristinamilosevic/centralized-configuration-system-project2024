package repositories

import (
	"errors"
	"fmt"
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

func (c ConfigInMemRepository) Add(config model.Config) {
	key := fmt.Sprintf("%s/%d", config.Name, config.Version)
	c.configs[key] = config
}

func (c ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	key := fmt.Sprintf("%s/%d", name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config{}, errors.New("config not found")
	}
	return config, nil
}

// GetAll vraća sve konfiguracije
func (repo *ConfigInMemRepository) GetAll() ([]model.Config, error) {
	configs := make([]model.Config, 0, len(repo.configs))
	for _, config := range repo.configs {
		configs = append(configs, config)
	}
	return configs, nil
}
