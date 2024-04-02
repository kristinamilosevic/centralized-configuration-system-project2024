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

func (s ConfigService) ReadConfigByName(name string) (model.Config, error) {
	return s.repo.ReadByName(name)
}

func (s ConfigService) UpdateConfig(config model.Config) error {
	return s.repo.Update(config)
}

func (s ConfigService) DeleteConfigByName(name string) error {
	return s.repo.DeleteByName(name)
}
