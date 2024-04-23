package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// Kanal za prekid signala
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Startovanje HTTP servera
	srv := &http.Server{Addr: ":8000"}

	repo := repositories.NewConfigInMemRepository()
	repoGroup := repositories.NewConfigGroupInMemRepository()
	service := services.NewConfigService(repo)
	serviceGroup := services.NewConfigGroupService(repoGroup)
	handler := handlers.NewConfigHandler(service)
	handlerGroup := handlers.NewConfigGroupHandler(serviceGroup)

	configs := []model.Config{}

	// Dodavanje pojedinačnih konfiguracija u listu
	params1 := map[string]string{"username": "pera", "password": "pera123"}
	config1 := model.Config{Name: "config1", Version: 1, Parameters: params1}
	configs = append(configs, config1)

	params2 := map[string]string{"username": "mika", "password": "mika123"}
	config2 := model.Config{Name: "config2", Version: 1, Parameters: params2}
	configs = append(configs, config2)

	params := map[string]string{"username": "pera", "password": "pera123"}
	config := model.Config{Name: "db_config", Version: 2, Parameters: params}

	service.Add(config)
	// Pravljenje konfiguracione grupe sa dodatom listom konfiguracija
	configGroup := model.ConfigGroup{Name: "configGroup", Version: 9, Configuration: configs}
	configGroup2 := model.ConfigGroup{Name: "configGroup2", Version: 2, Configuration: configs}

	serviceGroup.Add(configGroup)
	serviceGroup.Add(configGroup2)

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Get).Methods("GET")
	router.HandleFunc("/configs", handler.GetAll).Methods("GET")
	router.HandleFunc("/configGroups", handlerGroup.GetAll).Methods("GET")
	router.HandleFunc("/configs", handler.Create).Methods("POST")
	router.HandleFunc("/configGroups", handlerGroup.Create).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Delete).Methods("DELETE")
	router.HandleFunc("/configs/{name}/{version}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}/removeConfig/{configName}/{configVersion}", handlerGroup.RemoveConfig).Methods("DELETE")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}/addConfig", handlerGroup.AddConfig).Methods("PUT")

	// Pokretanje servera u zasebnoj gorutini
	go func() {
		log.Println("Starting server...")
		if err := http.ListenAndServe(":8000", router); err != nil {
			log.Fatal(err)
		}
	}()

	// Čekanje na prekid signala za graceful shutdown
	<-interrupt
	log.Println("Received SIGINT or SIGTERM. Shutting down...")

	// Shutdown servera
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}

	log.Println("Server successfully shut down.")
}
