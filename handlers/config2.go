package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"

	"github.com/gorilla/mux"
)

type Config2Handler struct {
	service services.Config2Service
}

func NewConfig2Handler(service services.Config2Service) Config2Handler {
	return Config2Handler{
		service: service,
	}
}

// swagger:route POST /configs configs createConfig
// Creates a new configuration.
//
// Responses:
//
//	201: NoContent
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
func (c Config2Handler) Create(w http.ResponseWriter, r *http.Request) {
	var config model.Config2
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Provera da li konfiguracija već postoji
	existingConfig, err := c.service.Get(config.Name, config.Version)
	if err == nil && (existingConfig.Name != "" || existingConfig.Version != 0) {
		http.Error(w, "configuration with this name and version already exists", http.StatusConflict)
		return
	}

	// Kreiranje nove konfiguracije
	err = c.service.CreateConfig(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// swagger:route GET /configs/{name}/{version} configs getConfig
// Get a configuration by name and version.
//
// Responses:
//
//	200: ResponseConfig2
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
func (c Config2Handler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// swagger:route DELETE /configs/{name}/{version} configs deleteConfig
// Deletes a configuration by name and version.
//
// Responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
func (c Config2Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Provera da li konfiguracija postoji pre nego što se pokuša brisanje
	_, err = c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	// Brisanje konfiguracije
	err = c.service.Delete(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /configs configs getAllConfigs
// Get all configurations.
//
// Responses:
//
//	200: []ResponseConfig2
//	500: InternalServerErrorResponse
func (c Config2Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	configs, err := c.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(configs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
