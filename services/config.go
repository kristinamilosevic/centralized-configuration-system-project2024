package services

import (
	"fmt"
	"projekat/model"
)

type ConfigService struct {
	repo model.ConfigRepository
}

func NewConfigService(repo model.ConfigRepository) ConfigService {
	return ConfigService{
		repo: repo,
	}
}

func (s ConfigService) Hello() {
	fmt.Println("hello from config service")
}

func (s ConfigService) CreateConfig(config model.Config) error {
	return s.repo.Create(config)
}

func (s ConfigService) Read(name string, version int) (model.Config, error) {
	return s.repo.Read(name, version)
}

func (s ConfigService) UpdateConfig(config model.Config) error {
	return s.repo.Update(config)
}

func (s ConfigService) Delete(name string, version int) error {
	return s.repo.Delete(name, version)
}

func (s ConfigService) Add(config model.Config) {
	s.repo.Add(config)
}

func (s ConfigService) Get(name string, version int) (model.Config, error) {
	return s.repo.Get(name, version)
}

func (s ConfigService) GetAll() ([]model.Config, error) {
	return s.repo.GetAll()
}
