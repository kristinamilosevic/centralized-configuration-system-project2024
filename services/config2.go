package services

import (
	"fmt"
	"projekat/model"
	"projekat/poststore"
)

type Config2Service struct {
	repo *poststore.PostStore
}

func NewConfig2Service(repo *poststore.PostStore) Config2Service {
	return Config2Service{
		repo: repo,
	}
}

func (s Config2Service) Hello() {
	fmt.Println("hello from config service")
}

func (s Config2Service) CreateConfig(config model.Config2, idempotencyKey, bodyHash string) error {
	return s.repo.CreateConfig(&config, idempotencyKey, bodyHash)
}

func (s Config2Service) GetByIdempotencyKey(idempotencyKey string) (string, error) {
	return s.repo.GetHashByIdempotencyKey(idempotencyKey)
}

func (s Config2Service) Read(name string, version int) (model.Config2, error) {
	config, err := s.repo.GetConfig(name, version)
	if err != nil {
		return model.Config2{}, err
	}
	return *config, nil
}
func (s Config2Service) CheckIfExists(idempotencyKey, bodyHash string) (bool, error) {
	return s.repo.CheckIfExists(idempotencyKey, bodyHash)
}
func (s Config2Service) UpdateConfig(config model.Config2) error {
	return s.repo.UpdateConfig(&config)
}

func (s Config2Service) Delete(name string, version int) error {
	_, err := s.repo.DeleteConfig(name, version)
	return err
}

func (s Config2Service) Get(name string, version int) (model.Config2, error) {
	config, err := s.repo.GetConfig(name, version)
	if err != nil {
		return model.Config2{}, err
	}
	return *config, nil
}

func (s Config2Service) GetAll() ([]model.Config2, error) {
	configs, err := s.repo.GetAllConfigs()
	if err != nil {
		return nil, err
	}
	var result []model.Config2
	for _, c := range configs {
		result = append(result, *c)
	}
	return result, nil
}
