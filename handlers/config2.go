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

// swagger:route POST /configs2 configs2 createConfig2
// Creates a new configuration.
//
// responses:
//
//	201: NoContent
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters createConfig2
type CreateConfig2Request struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestConfig"
	//  required: true
	Body model.Config2 `json:"body"`
}

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

// swagger:route GET /configs2/{name}/{version} configs2 getConfig2
// Get a configuration by name and version.
//
// responses:
//
//	200: ResponseConfig2
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters getConfig2
type GetConfig2Request struct {
	// Configuration name
	// in: path
	// required: true
	Name string `json:"name"`

	// Configuration version
	// in: path
	// required: true
	Version int `json:"version"`
}

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

// swagger:route DELETE /configs2/{name}/{version} configs2 deleteConfig2
// Deletes a configuration by name and version.
//
// responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters deleteConfig2
type DeleteConfig2Request struct {
	// Configuration name
	// in: path
	// required: true
	Name string `json:"name"`

	// Configuration version
	// in: path
	// required: true
	Version int `json:"version"`
}

func (c Config2Handler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the configuration exists before attempting deletion
	_, err = c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, "config not found", http.StatusNotFound)
		return
	}

	// Delete the configuration
	err = c.service.Delete(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /configs2 configs2 getAllConfigs
// Get all configurations.
//
// responses:
//
//	200: []ResponseConfig2
//	500: InternalServerErrorResponse
//
// swagger:parameters getAllConfigs2
type GetAllConfigs2Request struct {
	// No additional parameters needed
}

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

// swagger:response ResponseConfig2
type ResponseConfig2 struct {
	// Configuration
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestConfig"
	//  required: true
	Body model.Config2 `json:"body"`
}

// swagger:response BadRequestResponse
type BadRequestResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`

	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response NotFoundResponse
type NotFoundResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`

	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response NoContent
type NoContent struct{}
