package services

import (
	"fmt"
	"projekat/model"
)

type ConfigGroupService struct {
	repo model.ConfigGroupRepository
}

func NewConfigGroupService(repo model.ConfigGroupRepository) ConfigGroupService {
	return ConfigGroupService{
		repo: repo,
	}
}

func (s ConfigGroupService) Hello() {
	fmt.Println("hello from config group service")
}

func (s ConfigGroupService) Create(configGroup model.ConfigGroup) error {
	return s.repo.Create(configGroup)
}

func (s ConfigGroupService) Read(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Read(name, version)
}

func (s ConfigGroupService) Update(configGroup model.ConfigGroup) error {
	return s.repo.Update(configGroup)
}

func (s ConfigGroupService) Delete(name string, version int) error {
	return s.repo.Delete(name, version)
}

func (s ConfigGroupService) GetAll() ([]model.ConfigGroup, error) {
	return s.repo.GetAll()
}

func (s ConfigGroupService) Add(configGroup model.ConfigGroup) {
	s.repo.Add(configGroup)
}

func (s ConfigGroupService) Get(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}

func (s ConfigGroupService) RemoveConfig(groupName string, groupVersion int, configName string, configVersion int) error {
	// Prvo dohvatimo grupu konfiguracija
	configGroup, err := s.repo.Get(groupName, groupVersion)
	if err != nil {
		return err
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

	// Ažurirajmo grupu konfiguracija u repozitoriju
	err = s.repo.Update(configGroup)
	if err != nil {
		return err
	}

	return nil
}

func (s ConfigGroupService) AddConfigs(groupName string, groupVersion int, config model.Config2) error {
	// Prvo dohvatimo grupu konfiguracija
	configGroup, err := s.repo.Get(groupName, groupVersion)
	if err != nil {
		return err
	}

	// Dodajemo nove konfiguracije u grupu
	configGroup.Configuration = append(configGroup.Configuration, config)

	// Ažurirajmo grupu konfiguracija u repozitoriju
	err = s.repo.Update(configGroup)
	if err != nil {
		return err
	}

	return nil
}

func (s ConfigGroupService) GetFilteredConfigs(name string, version int, filter map[string]string) ([]model.Config2, error) {
	// Pozivamo odgovarajuću funkciju u repozitorijumu da bismo dobili filtrirane konfiguracije
	filteredConfigs, err := s.repo.GetFilteredConfigs(name, version, filter)
	if err != nil {
		return nil, err
	}
	return filteredConfigs, nil
}
