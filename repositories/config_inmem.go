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
	key := configKey(config.Name, config.Version)
	if _, exists := repo.configs[key]; exists {
		return errors.New("config with this name and version already exists")
	}

	repo.configs[key] = config
	return nil
}

func (repo *ConfigInMemRepository) Read(name string, version int) (model.Config, error) {
	for _, config := range repo.configs {
		if config.Name == name && config.Version == version {
			return config, nil
		}
	}
	return model.Config{}, errors.New("config not found")
}

func (repo *ConfigInMemRepository) Update(config model.Config) error {
	key := configKey(config.Name, config.Version)
	if _, exists := repo.configs[key]; !exists {
		return errors.New("config not found")
	}

	repo.configs[key] = config
	return nil
}

func (repo *ConfigInMemRepository) Delete(name string, version int) error {
	found := false
	for key, config := range repo.configs {
		if config.Name == name && config.Version == version {
			delete(repo.configs, key)
			found = true
		}
	}
	if !found {
		return errors.New("config not found")
	}
	return nil
}

func (c ConfigInMemRepository) Add(config model.Config) {
	key := configKey(config.Name, config.Version)
	c.configs[key] = config
}

func (c ConfigInMemRepository) Get(name string, version int) (model.Config, error) {
	key := configKey(name, version)
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

// configKey kreira ključ za konfiguraciju na osnovu imena i verzije
func configKey(name string, version int) string {
	return fmt.Sprintf("%s/%d", name, version)
}
