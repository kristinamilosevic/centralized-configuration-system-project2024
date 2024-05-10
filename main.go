package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"projekat/handlers"
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
	repo2 := repositories.NewConfig2InMemRepository()
	repoGroup := repositories.NewConfigGroupInMemRepository()
	service := services.NewConfigService(repo)
	service2 := services.NewConfig2Service(repo2)
	serviceGroup := services.NewConfigGroupService(repoGroup)
	handler := handlers.NewConfigHandler(service)
	handler2 := handlers.NewConfig2Handler(service2)
	handlerGroup := handlers.NewConfigGroupHandler(serviceGroup)

	router := mux.NewRouter()
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.HandleFunc("/configs2/{name}/{version}", handler2.Get).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Get).Methods("GET")
	router.HandleFunc("/configs", handler.GetAll).Methods("GET")
	router.HandleFunc("/configs2", handler2.GetAll).Methods("GET")
	router.HandleFunc("/configGroups", handlerGroup.GetAll).Methods("GET")
	router.HandleFunc("/configs", handler.Create).Methods("POST")
	router.HandleFunc("/configs2", handler2.Create).Methods("POST")
	router.HandleFunc("/configGroups", handlerGroup.Create).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Delete).Methods("DELETE")
	router.HandleFunc("/configs/{name}/{version}", handler.Delete).Methods("DELETE")
	router.HandleFunc("/configs2/{name}/{version}", handler2.Delete).Methods("DELETE")
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
