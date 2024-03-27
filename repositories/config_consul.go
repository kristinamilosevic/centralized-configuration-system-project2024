package repositories

import "projekat/model"

type ConfigConsulRepository struct {
}

//dodaj implementaciju metoda iz interfejsa ConfigRepo

func NewConfigConsulRepository() model.ConfigRepository {
	return ConfigConsulRepository{}
}
