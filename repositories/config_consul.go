package repositories

import (
	"errors"
	"projekat/model"
)

type ConfigConsulRepository struct {
}

//dodaj implementaciju metoda iz interfejsa ConfigRepo

func NewConfigConsulRepository() model.ConfigRepository {
	return &ConfigConsulRepository{}
}

func (repo *ConfigConsulRepository) Create(config model.Config) error {
	return errors.New("not implemented")
}

func (repo *ConfigConsulRepository) ReadByName(name string) (model.Config, error) {
	return model.Config{}, errors.New("not implemented")
}

func (repo *ConfigConsulRepository) Update(config model.Config) error {
	return errors.New("not implemented")
}

func (repo *ConfigConsulRepository) DeleteByName(name string) error {
	return errors.New("not implemented")
}
