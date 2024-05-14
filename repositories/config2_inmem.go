package repositories

import (
	"errors"
	"fmt"
	"projekat/model"
)

type Config2InMemRepository struct {
	configs map[string]model.Config2
}

func NewConfig2InMemRepository() model.Config2Repository {
	return &Config2InMemRepository{
		configs: make(map[string]model.Config2),
	}
}

func (repo *Config2InMemRepository) Create(config model.Config2) error {
	key := configKey2(config.Name, config.Version)
	if _, exists := repo.configs[key]; exists {
		return errors.New("config with this name and version already exists")
	}

	repo.configs[key] = config
	return nil
}

func (repo *Config2InMemRepository) Read(name string, version int) (model.Config2, error) {
	for _, config := range repo.configs {
		if config.Name == name && config.Version == version {
			return config, nil
		}
	}
	return model.Config2{}, errors.New("config not found")
}

func (repo *Config2InMemRepository) Update(config model.Config2) error {
	key := configKey2(config.Name, config.Version)
	if _, exists := repo.configs[key]; !exists {
		return errors.New("config not found")
	}

	repo.configs[key] = config
	return nil
}

func (repo *Config2InMemRepository) Delete(name string, version int) error {
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

func (c Config2InMemRepository) Add(config model.Config2) {
	key := configKey2(config.Name, config.Version)
	c.configs[key] = config
}

func (c Config2InMemRepository) Get(name string, version int) (model.Config2, error) {
	key := configKey2(name, version)
	config, ok := c.configs[key]
	if !ok {
		return model.Config2{}, errors.New("config not found")
	}
	return config, nil
}

// GetAll vraća sve konfiguracije
func (repo *Config2InMemRepository) GetAll() ([]model.Config2, error) {
	configs := make([]model.Config2, 0, len(repo.configs))
	for _, config := range repo.configs {
		configs = append(configs, config)
	}
	return configs, nil
}

// configKey kreira ključ za konfiguraciju na osnovu imena i verzije
func configKey2(name string, version int) string {
	return fmt.Sprintf("%s/%d", name, version)
}
