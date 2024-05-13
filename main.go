package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"projekat/handlers"
	"projekat/middleware"
	"projekat/model"
	"projekat/repositories"
	"projekat/services"

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
	repo2.Add(model.Config2{
		Name:    "config1",
		Version: 1,
	})
	repoGroup := repositories.NewConfigGroupInMemRepository()
	service := services.NewConfigService(repo)
	service2 := services.NewConfig2Service(repo2)
	serviceGroup := services.NewConfigGroupService(repoGroup)
	handler := handlers.NewConfigHandler(service)
	handler2 := handlers.NewConfig2Handler(service2)
	handlerGroup := handlers.NewConfigGroupHandler(serviceGroup)

	router := mux.NewRouter()

	// Kreiranje rate limiter middleware
	limiter := middleware.NewRateLimiter(5, time.Minute) // 5 zahteva u minuti

	// Dodavanje middleware-a na ruter
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter.ServeHTTP(w, r, next.ServeHTTP)
		})
	})

	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.HandleFunc("/configs2/{name}/{version}", handler2.Get).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Get).Methods("GET")
	router.HandleFunc("/configs", handler.GetAll).Methods("GET")
	router.HandleFunc("/configs2", handler2.GetAll).Methods("GET")
	router.HandleFunc("/configGroups", handlerGroup.GetAll).Methods("GET")
	router.HandleFunc("/configGroups/{name}/{version}/configs2/{filter}", handlerGroup.GetFilteredConfigs).Methods("GET")
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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server shutdown failed: %v", err)
	}

	log.Println("Server successfully shut down.")
}
