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
	"projekat/metrics"
	"projekat/middleware"
	"projekat/poststore"
	"projekat/repositories"
	"projekat/services"

	"github.com/gorilla/mux"
)

func main() {
	// Kanal za prekid signala
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

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

	// Initialize metrics
	metrics.Init()

	router := mux.NewRouter()

	// Kreiranje rate limiter middleware
	limiter := middleware.NewRateLimiter(100, time.Minute) // 5 zahteva u minuti

	// Dodavanje middleware-a na ruter
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			limiter.ServeHTTP(w, r, next.ServeHTTP)
		})
	})

	// Postavljanje ruta sa instrumentacijom
	router.HandleFunc("/configs/{name}/{version}", InstrumentHandler("GetConfig", handler.Get)).Methods("GET")
	router.HandleFunc("/configs", InstrumentHandler("GetAllConfigs", handler.GetAll)).Methods("GET")
	router.HandleFunc("/configs", InstrumentHandler("CreateConfig", handler.Create)).Methods("POST")
	router.HandleFunc("/configs/{name}/{version}", InstrumentHandler("DeleteConfig", handler.Delete)).Methods("DELETE")

	router.HandleFunc("/configs2/{name}/{version}", InstrumentHandler("GetConfig2", handler2.Get)).Methods("GET")
	router.HandleFunc("/configs2", InstrumentHandler("GetAllConfigs2", handler2.GetAll)).Methods("GET")
	router.HandleFunc("/configs2", InstrumentHandler("CreateConfig2", handler2.Create)).Methods("POST")
	router.HandleFunc("/configs2/{name}/{version}", InstrumentHandler("DeleteConfig2", handler2.Delete)).Methods("DELETE")

	router.HandleFunc("/configGroups/{name}/{version}", InstrumentHandler("GetConfigGroup", handlerGroup.Get)).Methods("GET")
	router.HandleFunc("/configGroups", InstrumentHandler("GetAllConfigGroups", handlerGroup.GetAll)).Methods("GET")
	router.HandleFunc("/configGroups", InstrumentHandler("CreateConfigGroup", handlerGroup.Create)).Methods("POST")
	router.HandleFunc("/configGroups/{name}/{version}", InstrumentHandler("DeleteConfigGroup", handlerGroup.Delete)).Methods("DELETE")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}/{configName}/{configVersion}", InstrumentHandler("RemoveConfigFromGroup", handlerGroup.RemoveConfig)).Methods("DELETE")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}", InstrumentHandler("AddConfigToGroup", handlerGroup.AddConfig)).Methods("PUT")
	router.HandleFunc("/configGroups/{name}/{version}/configs2/{filter}", InstrumentHandler("GetFilteredConfigs", handlerGroup.GetFilteredConfigs)).Methods("GET")
	router.HandleFunc("/configGroups/{groupName}/{groupVersion}/{filter}", InstrumentHandler("RemoveConfigsByLabels", handlerGroup.RemoveByLabels)).Methods("DELETE")

	// Metrics endpoint
	router.Handle("/metrics", metrics.MetricsHandler())

	// Startovanje HTTP servera
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

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

// InstrumentHandler instruments an HTTP handler with Prometheus metrics
func InstrumentHandler(name string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Mera vremena početka obrade zahteva

		// Instrumentovanje ukupnog broja zahteva
		metrics.RequestTotal.WithLabelValues(r.Method, name).Inc()

		// Hvatanje statusnog koda
		rr := &responseRecorder{w, http.StatusOK}
		h.ServeHTTP(rr, r)

		duration := time.Since(start).Seconds()                                   // Izračunavanje trajanja obrade zahteva u sekundama
		metrics.RequestDuration.WithLabelValues(r.Method, name).Observe(duration) // Merenje vremena odziva zahteva

		// Kategorizacija zahteva kao uspešan ili neuspešan
		if rr.statusCode >= 200 && rr.statusCode < 400 {
			metrics.RequestSuccessTotal.WithLabelValues(r.Method, name).Inc()
		} else {
			metrics.RequestFailureTotal.WithLabelValues(r.Method, name).Inc()
		}

		// Izračunavanje broja zahteva u sekundi (za jednostavnost, inkrementuje se za 1/duration)
		metrics.RequestsPerSecond.WithLabelValues(r.Method, name).Add(1 / duration)
	}
}

// responseRecorder is a wrapper to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
