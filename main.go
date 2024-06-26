// Post API
//
//	Title: Post API
//
//	Schemes: http
//	Version: 0.0.1
//	BasePath: /
//
//	Produces:
//	  - application/json
//
// swagger:meta
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
	"projekat/poststore"
	"projekat/rate_limiter"
	"projekat/repositories"
	"projekat/services"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func main() {
	// Kanal za prekid signala
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	// Startovanje HTTP servera
	srv := &http.Server{Addr: ":8000"}

	// Inicijalizacija poststore-a
	store, err := poststore.New()
	if err != nil {
		log.Fatalf("Failed to initialize poststore: %v", err)
	}
	storeGroup, err := poststore.NewGroupStore()
	if err != nil {
		log.Fatalf("Failed to initialize poststore: %v", err)
	}

	// Inicijalizacija in-memory skladišta
	repo := repositories.NewConfigInMemRepository()
	//repoGroup := repositories.NewConfigGroupInMemRepository()

	// Inicijalizacija servisa
	service := services.NewConfigService(repo)
	service2 := services.NewConfig2Service(store)
	serviceGroup := services.NewConfigGroupService(storeGroup)

	// Inicijalizacija handlera
	handler := handlers.NewConfigHandler(service)
	handler2 := handlers.NewConfig2Handler(service2)
	handlerGroup := handlers.NewConfigGroupHandler(serviceGroup)

	router := mux.NewRouter()

	// Kreiranje rate limiter middleware
	limiter := rate_limiter.NewRateLimiter(100, time.Minute) // 5 zahteva u minuti

	// Dodavanje middleware-a na ruter
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter.ServeHTTP(w, r, next.ServeHTTP)
		})
	})

	// Postavljanje svih ruta
	router.HandleFunc("/configs/{name}/{version}", handler.Get).Methods("GET")
	router.HandleFunc("/configs", handler.GetAll).Methods("GET")
	router.HandleFunc("/configs", handler.Create).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", handler.Delete).Methods("DELETE")

	router.HandleFunc("/configs2/{name}/{version}", handler2.Get).Methods("GET")
	router.HandleFunc("/configs2", handler2.GetAll).Methods("GET")
	router.HandleFunc("/configs2", handler2.Create).Methods("POST")
	router.HandleFunc("/configs2/{name}/{version}", handler2.Delete).Methods("DELETE")

	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Get).Methods("GET")
	router.HandleFunc("/configGroups", handlerGroup.GetAll).Methods("GET")
	router.HandleFunc("/configGroups", handlerGroup.Create).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", handlerGroup.Delete).Methods("DELETE")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}/{configName}/{configVersion}", handlerGroup.RemoveConfig).Methods("DELETE")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}", handlerGroup.AddConfig).Methods("PUT")
	router.HandleFunc("/configGroups/{name}/{version}/configs2/{filter}", handlerGroup.GetFilteredConfigs).Methods("GET")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}/{filter}", handlerGroup.RemoveByLabels).Methods("DELETE")

	// SwaggerUI
	optionsDevelopers := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	developerDocumentationHandler := middleware.SwaggerUI(optionsDevelopers, nil)
	router.Handle("/docs", developerDocumentationHandler)

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
