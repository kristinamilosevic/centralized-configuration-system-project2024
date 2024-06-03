package services

import (
	"errors"
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

	// Provera da li konfiguracija već postoji
	for _, existingConfig := range configGroup.Configuration {
		if existingConfig.Name == config.Name && existingConfig.Version == config.Version {
			return errors.New("configuration with the same name and version already exists")
		}
	}

	// Dodajemo novu konfiguraciju u grupu
	configGroup.Configuration = append(configGroup.Configuration, config)

	// Ažuriramo grupu konfiguracija u repozitoriju
	err = s.repo.Update(configGroup)
	if err != nil {
		return err
	}

	return nil
}

func (s ConfigGroupService) GetFilteredConfigs(name string, version int, filter map[string]string) ([]model.Config2, error) {
	filteredConfigs, err := s.repo.GetFilteredConfigs(name, version, filter)
	if err != nil {
		return nil, err
	}
	return filteredConfigs, nil
}

func (s ConfigGroupService) RemoveByLabels(groupName string, groupVersion int, filter map[string]string) error {
	// Dohvatimo konfiguracionu grupu
	configGroup, err := s.repo.Get(groupName, groupVersion)
	if err != nil {
		return err
	}

	// Inicijalizujemo slice za čuvanje preostalih konfiguracija
	var remainingConfigs []model.Config2

	// Iteriramo kroz sve konfiguracije i proveravamo da li odgovaraju filteru
	for _, config := range configGroup.Configuration {
		if !labelsExactMatch(config.Labels, filter) {
			remainingConfigs = append(remainingConfigs, config)
		}
	}

	// Ako sve konfiguracije odgovaraju filteru, vratimo odgovarajuću grešku
	if len(remainingConfigs) == len(configGroup.Configuration) {
		return errors.New("no configurations found matching the provided labels")
	}

	// Ažuriramo grupu konfiguracija sa preostalim konfiguracijama
	configGroup.Configuration = remainingConfigs

	// Ažuriramo grupu konfiguracija u repozitoriju
	err = s.repo.Update(configGroup)
	if err != nil {
		return fmt.Errorf("failed to update config group: %v", err)
	}

	return nil
}

func labelsExactMatch(configLabels, filter map[string]string) bool {
	if len(configLabels) != len(filter) {
		return false
	}
	for key, value := range filter {
		if configLabels[key] != value {
			return false
		}
	}
	return true
}
