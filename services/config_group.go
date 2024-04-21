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
