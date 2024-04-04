package main

import (
	"fmt"
	"net/http"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"

	"github.com/gorilla/mux"
)

func main() {
	repo := repositories.NewConfigInMemRepository()
	service := services.NewConfigService(repo)
	handler := handlers.NewConfigHandler(service)
	params := make(map[string]string)
	params["username"] = "pera"
	params["password"] = "pera123"
	config := model.Config{
		Name:       "db_config",
		Version:    2,
		Parameters: params,
	}

	service.Add(config)

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")

	http.ListenAndServe("0.0.0.0:8000", router)

	configGroupRepo := repositories.NewConfigGroupInMemRepository()
	configGroupService := services.NewConfigGroupService(configGroupRepo)

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

	configGroupService.CreateConfigGroup(model.ConfigGroup{Name: "FirstConfigurationGroup", Version: 1, Configuration: []model.Config{}})
	configGroup, err := configGroupService.ReadConfigGroupByName("FirstConfigurationGroup")
	if err != nil {
		fmt.Println("Error while reading configuration group:", err)
	} else {
		fmt.Println("Read configuration group:", configGroup)
	}

	configGroupService.CreateConfigGroup(model.ConfigGroup{Name: "SecondConfigurationGroup", Version: 1, Configuration: []model.Config{}})
	configGroup2, err := configGroupService.ReadConfigGroupByName("SecondConfigurationGroup")
	if err != nil {
		fmt.Println("Error while reading configuration group:", err)
	} else {
		fmt.Println("Read configuration group:", configGroup2)
	}

	configGroupService.UpdateConfigGroup(model.ConfigGroup{Name: "FirstConfigurationGroup", Version: 2, Configuration: []model.Config{}})
	err = configGroupService.DeleteConfigGroupByName("FirstConfigurationGroup")
	if err != nil {
		fmt.Println("Error while deleting configuration group:", err)
	} else {
		fmt.Println("Configuration group successfully deleted")
	}
}
