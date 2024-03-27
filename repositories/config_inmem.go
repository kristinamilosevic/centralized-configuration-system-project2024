package repositories

import "projekat/model"

type ConfigInMemRepository struct {
	//lista ili mapa nasih konfiguracija
	configs map[string]model.Config
}

//dodaj implementaciju metoda iz interfejsa

func NewConfigInMemRepository() model.ConfigRepository {
	return ConfigInMemRepository{
		configs: make(map[string]model.Config),
	}
}
