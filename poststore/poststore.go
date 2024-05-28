package poststore

import (
	"encoding/json"
	"fmt"
	"os"

	"projekat/model"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

type PostStore struct {
	cli *api.Client
}

func New() (*PostStore, error) {
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
	return &PostStore{cli: client}, nil
}

func generateKey() (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf("configs/%s", id), id
}

func constructKey(name string, version int) string {
	return fmt.Sprintf("configs/%s/%d", name, version)
}

func (ps *PostStore) CreateConfig(config *model.Config2) (*model.Config2, error) {
	kv := ps.cli.KV()
	key := constructKey(config.Name, config.Version)
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (ps *PostStore) GetConfig(name string, version int) (*model.Config2, error) {
	kv := ps.cli.KV()
	pair, _, err := kv.Get(constructKey(name, version), nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, fmt.Errorf("config not found")
	}
	var config model.Config2
	err = json.Unmarshal(pair.Value, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (ps *PostStore) GetAllConfigs() ([]*model.Config2, error) {
	kv := ps.cli.KV()
	data, _, err := kv.List("configs/", nil)
	if err != nil {
		return nil, err
	}
	var configs []*model.Config2
	for _, pair := range data {
		var config model.Config2
		err = json.Unmarshal(pair.Value, &config)
		if err != nil {
			continue
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

func (ps *PostStore) UpdateConfig(config *model.Config2) error {
	kv := ps.cli.KV()
	key := constructKey(config.Name, config.Version)
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	p := &api.KVPair{Key: key, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PostStore) DeleteConfig(name string, version int) (map[string]string, error) {
	kv := ps.cli.KV()
	_, err := kv.Delete(constructKey(name, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"Deleted": fmt.Sprintf("%s/%d", name, version)}, nil
}
