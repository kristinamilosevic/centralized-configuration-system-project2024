package services

import (
	"fmt"
	"projekat/model"
)

type Config2Service struct {
	repo model.Config2Repository
}

func NewConfig2Service(repo model.Config2Repository) Config2Service {
	return Config2Service{
		repo: repo,
	}
}

func (s Config2Service) Hello() {
	fmt.Println("hello from config service")
}

func (s Config2Service) CreateConfig(config model.Config2) error {
	return s.repo.Create(config)
}

func (s Config2Service) Read(name string, version int) (model.Config2, error) {
	return s.repo.Read(name, version)
}

func (s Config2Service) UpdateConfig(config model.Config2) error {
	return s.repo.Update(config)
}

func (s Config2Service) Delete(name string, version int) error {
	return s.repo.Delete(name, version)
}

func (s Config2Service) Add(config model.Config2) {
	s.repo.Add(config)
}

func (s Config2Service) Get(name string, version int) (model.Config2, error) {
	return s.repo.Get(name, version)
}

func (s Config2Service) GetAll() ([]model.Config2, error) {
	return s.repo.GetAll()
}
