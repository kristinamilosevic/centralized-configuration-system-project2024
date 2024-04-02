package main

import (
	"fmt"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)

	service.Hello()

	service.CreateConfig(model.Config{Name: "FirstConfiguration", Version: 1, Parameters: map[string]string{"parameter": "value"}})
	config, err := service.ReadConfigByName("FirstConfiguration")
	if err != nil {
		fmt.Println("Error while reading configuration:", err)
	} else {
		fmt.Println("Read configuration:", config)
	}

	service.CreateConfig(model.Config{Name: "SecondConfiguration", Version: 1, Parameters: map[string]string{"parameter": "value"}})
	config2, err := service.ReadConfigByName("SecondConfiguration")
	if err != nil {
		fmt.Println("Error while reading configuration:", err)
	} else {
		fmt.Println("Read configuration:", config2)
	}

	service.UpdateConfig(model.Config{Name: "FirstConfiguration", Version: 2, Parameters: map[string]string{"parameter": "new_value"}})

	err = service.DeleteConfigByName("FirstConfiguration")
	if err != nil {
		fmt.Println("Error while deleting configuration:", err)
	} else {
		fmt.Println("Configuration successfully deleted")
	}
}
