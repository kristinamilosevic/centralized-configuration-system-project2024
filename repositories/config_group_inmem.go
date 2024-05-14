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
	key := configGroupKey(configGroup.Name, configGroup.Version)
	if _, exists := repo.configGroups[key]; exists {
		return errors.New("config group with this name already exists")
	}

	repo.configGroups[key] = configGroup
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
	found := false
	for key, configGroup := range repo.configGroups {
		if configGroup.Name == name && configGroup.Version == version {
			delete(repo.configGroups, key)
			found = true
		}
	}
	if !found {
		return errors.New("config group not found")
	}
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

func (repo *ConfigGroupInMemRepository) RemoveConfig(groupName string, groupVersion int, configName string, configVersion int) error {
	key := configGroupKey(groupName, groupVersion)
	configGroup, ok := repo.configGroups[key]
	if !ok {
		return errors.New("config group not found")
	}

	// Pronađimo konfiguraciju koju želimo ukloniti iz grupe
	var indexToRemove = -1
	for i, config := range configGroup.Configuration {
		if config.Name == configName && config.Version == configVersion {
			indexToRemove = i
			break
		}
	}
	if indexToRemove == -1 {
		return fmt.Errorf("config with name %s and version %d not found in group", configName, configVersion)
	}

	// Uklonimo konfiguraciju iz grupe
	configGroup.Configuration = append(configGroup.Configuration[:indexToRemove], configGroup.Configuration[indexToRemove+1:]...)

	// Ažurirajmo grupu konfiguracija u memoriji
	repo.configGroups[key] = configGroup

	return nil
}

func (repo *ConfigGroupInMemRepository) AddConfig(groupName string, groupVersion int, config model.Config2) error {
	// Kreiramo ključ za grupu konfiguracija
	key := configGroupKey(groupName, groupVersion)

	// Proveravamo postoji li već grupa sa tim ključem
	configGroup, ok := repo.configGroups[key]
	if !ok {
		return errors.New("config group not found")
	}

	// Kreiramo novu konfiguraciju
	newConfig := model.NewConfig2_2(groupName, groupVersion) // Prilagoditi kreiranje nove konfiguracije prema potrebama

	// Dodajemo novu konfiguraciju u grupu
	configGroup.Configuration = append(configGroup.Configuration, newConfig)

	// Ažuriramo grupu konfiguracija u memoriji
	repo.configGroups[key] = configGroup

	return nil
}

func (repo *ConfigGroupInMemRepository) GetFilteredConfigs(name string, version int, filter map[string]string) ([]model.Config2, error) {
	// Dobijanje konfiguracione grupe
	configGroup, err := repo.Read(name, version)
	if err != nil {
		return nil, err
	}

	filteredConfigs := make([]model.Config2, 0)

	// Iteriranje kroz sve konfiguracije u grupi i provera da li se poklapaju sa filterom
	for _, config := range configGroup.Configuration {
		matches := true
		for key, value := range filter {
			if config.Labels[key] != value {
				matches = false
				break
			}
		}
		if matches {
			filteredConfigs = append(filteredConfigs, config)
		}
	}

	return filteredConfigs, nil
}

func (repo *ConfigGroupInMemRepository) RemoveByLabels(groupName string, groupVersion int, filter map[string]string) error {
	key := configGroupKey(groupName, groupVersion)
	configGroup, ok := repo.configGroups[key]
	if !ok {
		return errors.New("config group not found")
	}

	// Inicijalizacija promenljive za praćenje broja uklonjenih konfiguracija
	removed := 0

	// Iteriranje kroz sve konfiguracije i provera da li odgovaraju filteru
	for i := len(configGroup.Configuration) - 1; i >= 0; i-- {
		config := configGroup.Configuration[i]
		matches := true
		for key, value := range filter {
			if config.Labels[key] != value {
				matches = false
				break
			}
		}
		if matches {
			// Uklanjanje konfiguracije
			configGroup.Configuration = append(configGroup.Configuration[:i], configGroup.Configuration[i+1:]...)
			removed++
		}
	}

	if removed == 0 {
		return errors.New("no configurations found matching the provided labels")
	}

	repo.configGroups[key] = configGroup
	return nil
}
