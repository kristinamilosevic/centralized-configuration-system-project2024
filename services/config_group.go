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

	// Inicijalizujemo slice za čuvanje indeksa elemenata koji treba ukloniti
	indicesToRemove := []int{}

	// Iteriramo kroz sve konfiguracije i proveravamo da li odgovaraju filteru
	for i, config := range configGroup.Configuration {
		matches := true
		for key, value := range filter {
			if config.Labels[key] != value {
				matches = false
				break
			}
		}
		if matches {
			indicesToRemove = append(indicesToRemove, i)
		}
	}

	// Ako nema pronađenih konfiguracija koje odgovaraju filteru, vratimo odgovarajuću grešku
	if len(indicesToRemove) == 0 {
		return errors.New("no configurations found matching the provided labels")
	}

	// Uklanjamo konfiguracije sa odgovarajućim indeksima iz konfiguracione grupe
	for i := len(indicesToRemove) - 1; i >= 0; i-- {
		index := indicesToRemove[i]
		configGroup.Configuration = append(configGroup.Configuration[:index], configGroup.Configuration[index+1:]...)
	}

	// Ažuriramo konfiguracionu grupu u repozitorijumu
	err = s.repo.Update(configGroup)
	if err != nil {
		return fmt.Errorf("failed to update config group: %v", err)
	}

	return nil
}
