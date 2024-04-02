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

func (s ConfigGroupService) CreateConfigGroup(configGroup model.ConfigGroup) error {
	return s.repo.Create(configGroup)
}

func (s ConfigGroupService) ReadConfigGroupByName(name string) (model.ConfigGroup, error) {
	return s.repo.ReadByName(name)
}

func (s ConfigGroupService) UpdateConfigGroup(configGroup model.ConfigGroup) error {
	return s.repo.Update(configGroup)
}

func (s ConfigGroupService) DeleteConfigGroupByName(name string) error {
	return s.repo.DeleteByName(name)
}
