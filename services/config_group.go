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

func (s ConfigGroupService) AddConfigToGroup(groupName string, config model.Config) error {
	group, err := s.repo.ReadByName(groupName)
	if err != nil {
		return err
	}

	group.Configuration = append(group.Configuration, config)

	err = s.repo.Update(group)
	if err != nil {
		return err
	}

	return nil
}

func (s ConfigGroupService) Get(name string, version int) (model.ConfigGroup, error) {
	return s.repo.Get(name, version)
}
func (s ConfigGroupService) GetAll() ([]model.ConfigGroup, error) {
	return s.repo.GetAll()
}

func (s ConfigGroupService) Add(configGroup model.ConfigGroup) {
	s.repo.Add(configGroup)
}
