package handlers

import (
	"crypto/sha256"
	"encoding/hex"
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

// Hash function for the request body
func hashRequestBody(body interface{}) (string, error) {
	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(jsonBytes)
	return hex.EncodeToString(hash[:]), nil
}

// POST /configs2
func (c Config2Handler) Create(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header is required", http.StatusBadRequest)
		return
	}

	var config model.Config2
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate hash of the request body
	bodyHash, err := hashRequestBody(config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the idempotency key and body hash combination already exists
	exists, err := c.service.CheckIfExists(idempotencyKey, bodyHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "Request with the same Idempotency-Key and body already exists", http.StatusConflict)
		return
	}

	// Create new configuration with the combination of Idempotency-Key and body hash
	err = c.service.CreateConfig(config, idempotencyKey, bodyHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GET /configs/{name}/{version}
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

// DELETE /configs/{name}/{version}
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

// GET /configs
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
