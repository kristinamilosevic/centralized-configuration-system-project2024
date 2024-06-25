package poststore

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"projekat/model"

	"github.com/hashicorp/consul/api"
)

type GroupStore struct {
	cli *api.Client
}

func NewGroupStore() (*GroupStore, error) {
	db := os.Getenv("DB")
	if db == "" {
		db = "localhost"
	}
	dbport := os.Getenv("DBPORT")
	if dbport == "" {
		dbport = "8500"
	}
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &GroupStore{cli: client}, nil
}

func constructGroupKey(name string, version int) string {
	return fmt.Sprintf("configGroups/%s/%d", name, version)
}

func (gs *GroupStore) CreateConfigGroup(configGroup *model.ConfigGroup, idempotencyKey, bodyHash string) error {
	kv := gs.cli.KV()
	key := constructGroupKey(configGroup.Name, configGroup.Version)
	idempotentKey := fmt.Sprintf("idempotency/%s/%s", idempotencyKey, bodyHash)

	data, err := json.Marshal(configGroup)
	if err != nil {
		return err
	}

	// Skladištenje konfiguracione grupe
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}

	// Skladištenje idempotentnog ključa i hash-a tela kao JSON objekat
	idempotentData, err := json.Marshal(map[string]string{"body-hash": bodyHash})
	if err != nil {
		return err
	}
	idempotentKV := &api.KVPair{Key: idempotentKey, Value: idempotentData}
	_, err = kv.Put(idempotentKV, nil)
	if err != nil {
		return err
	}

	return nil
}

func (gs *GroupStore) CheckIfExists(idempotencyKey, bodyHash string) (bool, error) {
	kv := gs.cli.KV()
	idempotentKey := fmt.Sprintf("idempotency/%s/%s", idempotencyKey, bodyHash)

	pair, _, err := kv.Get(idempotentKey, nil)
	if err != nil {
		return false, err
	}
	if pair == nil {
		return false, nil
	}

	return true, nil
}

func (gs *GroupStore) Read(name string, version int) (model.ConfigGroup, error) {
	kv := gs.cli.KV()
	key := constructGroupKey(name, version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return model.ConfigGroup{}, err
	}
	if pair == nil {
		return model.ConfigGroup{}, errors.New("config group not found")
	}
	var configGroup model.ConfigGroup
	err = json.Unmarshal(pair.Value, &configGroup)
	return configGroup, err
}

func (gs *GroupStore) Update(newConfigGroup model.ConfigGroup) error {
	kv := gs.cli.KV()
	key := constructGroupKey(newConfigGroup.Name, newConfigGroup.Version)
	pair, _, err := kv.Get(key, nil)
	if err != nil {
		return err
	}
	if pair == nil {
		return errors.New("config group not found")
	}
	data, err := json.Marshal(newConfigGroup)
	if err != nil {
		return err
	}
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	return err
}

func (gs *GroupStore) Delete(name string, version int) error {
	kv := gs.cli.KV()
	key := constructGroupKey(name, version)
	_, err := kv.Delete(key, nil)
	return err
}

func (gs *GroupStore) GetAll() ([]model.ConfigGroup, error) {
	kv := gs.cli.KV()
	data, _, err := kv.List("configGroups/", nil)
	if err != nil {
		return nil, err
	}
	var configGroups []model.ConfigGroup
	for _, pair := range data {
		var configGroup model.ConfigGroup
		err = json.Unmarshal(pair.Value, &configGroup)
		if err != nil {
			continue
		}
		configGroups = append(configGroups, configGroup)
	}
	return configGroups, nil
}

func (gs *GroupStore) Add(configGroup model.ConfigGroup) {
	// Add method can be implemented if needed, similar to Create
}

func (gs *GroupStore) Get(name string, version int) (model.ConfigGroup, error) {
	return gs.Read(name, version)
}

func (gs *GroupStore) RemoveConfig(groupName string, groupVersion int, configName string, configVersion int) error {
	configGroup, err := gs.Read(groupName, groupVersion)
	if err != nil {
		return err
	}
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
	configGroup.Configuration = append(configGroup.Configuration[:indexToRemove], configGroup.Configuration[indexToRemove+1:]...)
	return gs.Update(configGroup)
}

func (gs *GroupStore) AddConfig(groupName string, groupVersion int, config model.Config2) error {
	configGroup, err := gs.Read(groupName, groupVersion)
	if err != nil {
		return err
	}

	// Provera da li konfiguracija već postoji
	for _, existingConfig := range configGroup.Configuration {
		if existingConfig.Name == config.Name && existingConfig.Version == config.Version {
			return errors.New("configuration with the same name and version already exists")
		}
	}

	// Dodavanje nove konfiguracije
	configGroup.Configuration = append(configGroup.Configuration, config)
	return gs.Update(configGroup)
}

func (gs *GroupStore) GetFilteredConfigs(name string, version int, filter map[string]string) ([]model.Config2, error) {
	// Dobavljanje konfiguracione grupe iz skladišta
	configGroup, err := gs.Read(name, version)
	if err != nil {
		return nil, err
	}

	var filteredConfigs []model.Config2

	// Kreiranje prefiksa na osnovu filtera
	prefix := ""
	for key, value := range filter {
		prefix = fmt.Sprintf("%s%s=%s", prefix, key, value)
	}

	// Filtriranje konfiguracija
	for _, config := range configGroup.Configuration {
		if len(config.Labels) != len(filter) {
			continue // preskoči ako broj etiketa ne odgovara
		}

		configPrefix := ""
		for key, value := range config.Labels {
			configPrefix = fmt.Sprintf("%s%s=%s", configPrefix, key, value)
		}

		if configPrefix == prefix {
			filteredConfigs = append(filteredConfigs, config)
		}
	}

	return filteredConfigs, nil
}

func (gs *GroupStore) RemoveByLabels(groupName string, groupVersion int, filter map[string]string) error {
	fmt.Println("RemoveByLabels method called with groupName:", groupName, "groupVersion:", groupVersion, "filter:", filter)

	// Čitanje grupe konfiguracija iz baze podataka
	configGroup, err := gs.Read(groupName, groupVersion)
	if err != nil {
		fmt.Println("Error reading config group:", err)
		return err
	}

	// Kreiranje novog niza za konfiguracije koje ostaju
	var remainingConfigs []model.Config2
	for _, config := range configGroup.Configuration {
		// Provera da li konfiguracija zadovoljava filter
		matches := true
		for key, value := range filter {
			if config.Labels[key] != value {
				matches = false
				break
			}
		}
		// Ako ne zadovoljava filter, dodajemo je u novi niz
		if !matches {
			remainingConfigs = append(remainingConfigs, config)
		}
	}

	// Ažuriranje konfiguracija u grupi samo sa onima koje nisu uklonjene
	configGroup.Configuration = remainingConfigs

	// Ako su sve konfiguracije uklonjene, vratimo grešku
	if len(configGroup.Configuration) == 0 {
		fmt.Println("No configurations found matching the provided labels")
		return errors.New("no configurations found matching the provided labels")
	}

	// Ažuriranje grupe konfiguracija u bazi podataka
	if err := gs.Update(configGroup); err != nil {
		fmt.Println("Error updating config group:", err)
		return err
	}

	fmt.Println("Config group successfully updated after removing configurations by labels")
	return nil
}
