package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type ConfigGroupHandler struct {
	service services.ConfigGroupService
}

func NewConfigGroupHandler(service services.ConfigGroupService) ConfigGroupHandler {
	return ConfigGroupHandler{
		service: service,
	}
}

// swagger:route POST /configGroups configGroups createConfigGroup
// Creates a new configuration group.
//
// responses:
//
//	201: NoContent
//	400: BadRequestResponse
//	409: ErrorResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters createConfigGroup
type CreateConfigGroupRequest struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestConfigGroup"
	//  required: true
	Body model.ConfigGroup `json:"body"`
}

func (c ConfigGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header is required", http.StatusBadRequest)
		return
	}

	var configGroup model.ConfigGroup
	err := json.NewDecoder(r.Body).Decode(&configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate hash of the request body
	bodyHash, err := hashRequestBody(configGroup)
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

	// Create new configuration group with the combination of Idempotency-Key and body hash
	err = c.service.CreateConfigGroup(configGroup, idempotencyKey, bodyHash)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// swagger:route GET /configGroups/{name}/{version} configGroups getConfigGroup
// Retrieves a configuration group by name and version.
//
// responses:
//
//	200: ResponseConfigGroup
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters getConfigGroup
type GetConfigGroupRequest struct {
	// Configuration group name
	// in: path
	// required: true
	Name string `json:"name"`

	// Configuration group version
	// in: path
	// required: true
	Version int `json:"version"`
}

func (c ConfigGroupHandler) Get(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	configGroup, err := c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// swagger:route DELETE /configGroups/{name}/{version} configGroups deleteConfigGroup
// Deletes a configuration group by name and version.
//
// responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters deleteConfigGroup
type DeleteConfigGroupRequest struct {
	// Configuration group name
	// in: path
	// required: true
	Name string `json:"name"`

	// Configuration group version
	// in: path
	// required: true
	Version int `json:"version"`
}

func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the configuration group exists before attempting deletion
	_, err = c.service.Get(name, versionInt)
	if err != nil {
		http.Error(w, "config group not found", http.StatusNotFound)
		return
	}

	// Delete the configuration group
	err = c.service.Delete(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /configGroups configGroups getAllConfigGroups
// Retrieves all configuration groups.
//
// responses:
//
//   200: getAllConfigGroupsResponse
//   500: InternalServerErrorResponse

// swagger:response getAllConfigGroupsResponse
type getAllConfigGroupsResponse struct {
	// Configuration groups
	// in: body
	Body []model.ConfigGroup `json:"body"`
}

func (c ConfigGroupHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	configGroups, err := c.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(configGroups)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// swagger:route DELETE /configGroups/{groupName}/{groupVersion}/{configName}/{configVersion} configGroups removeConfigFromGroup
// Removes a configuration from a group by name and version.
//
// responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters removeConfigFromGroup
type RemoveConfigFromGroupRequest struct {
	// Configuration group name
	// in: path
	// required: true
	GroupName string `json:"groupName"`

	// Configuration group version
	// in: path
	// required: true
	GroupVersion string `json:"groupVersion"`

	// Configuration name
	// in: path
	// required: true
	ConfigName string `json:"configName"`

	// Configuration version
	// in: path
	// required: true
	ConfigVersion string `json:"configVersion"`
}

func (c ConfigGroupHandler) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["groupName"]
	groupVersion := mux.Vars(r)["groupVersion"]
	configName := mux.Vars(r)["configName"]
	configVersion := mux.Vars(r)["configVersion"]

	groupVersionInt, err := strconv.Atoi(groupVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	configVersionInt, err := strconv.Atoi(configVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.RemoveConfig(groupName, groupVersionInt, configName, configVersionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route PUT /configGroups/{groupName}/{groupVersion} configGroups addConfigToGroup
// Adds a configuration to a group by name and version.
//
// responses:
//
//	201: NoContent
//	400: BadRequestResponse
//	409: ErrorResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters addConfigToGroup
type AddConfigToGroupRequest struct {
	// Configuration group name
	// in: path
	GroupName string `json:"groupName"`

	// Configuration group version
	// in: path
	GroupVersion string `json:"groupVersion"`

	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/RequestConfigGroup"
	//  required: true
	Body model.Config2 `json:"body"`
}

func (c ConfigGroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["groupName"]
	groupVersion := mux.Vars(r)["groupVersion"]

	groupVersionInt, err := strconv.Atoi(groupVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config := model.Config2{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	configGroup, err := c.service.Get(groupName, groupVersionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, existingConfig := range configGroup.Configuration {
		if existingConfig.Name == config.Name && existingConfig.Version == config.Version {
			http.Error(w, "configuration with the same name and version already exists", http.StatusConflict)
			return
		}
	}

	err = c.service.AddConfigs(groupName, groupVersionInt, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// swagger:route GET /configGroups/{name}/{version}/{filter} configGroups getFilteredConfigs
// Retrieves filtered configurations from a group by name, version, and filter criteria.
//
// responses:
//
//	200: []ResponseConfig
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters getFilteredConfigs
type GetFilteredConfigsRequest struct {
	// Ime grupe konfiguracija
	// in: path
	// required: true
	Name string `json:"name"`

	// Verzija grupe konfiguracija
	// in: path
	// required: true
	Version string `json:"version"`

	// Filter za konfiguracije
	// in: path
	// required: true
	Filter string `json:"filter"`
}

// swagger:response ResponseConfig
type ResponseConfig struct {
	// Id konfiguracije
	// in: string
	ID string `json:"id"`

	// Naziv konfiguracije
	// in: string
	Name string `json:"name"`

	// Verzija konfiguracije
	// in: int
	Version int `json:"version"`

	// Opis konfiguracije
	// in: string
	Description string `json:"description"`

	// Lista tagova konfiguracije
	// in: []string
	Tags []string `json:"tags"`
}

func (c ConfigGroupHandler) GetFilteredConfigs(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	versionStr := mux.Vars(r)["version"]
	filterStr := mux.Vars(r)["filter"]

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterMap := make(map[string]string)
	if filterStr != "" {
		keyValues := strings.Split(filterStr, ",")
		for _, kv := range keyValues {
			parts := strings.Split(kv, "=")
			if len(parts) == 2 {
				filterMap[parts[0]] = parts[1]
			}
		}
	}

	configs, err := c.service.GetFilteredConfigs(name, version, filterMap)
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

// swagger:route DELETE /configGroups/{groupName}/{groupVersion}/{filter} configGroups removeByLabels
// Removes configurations from a group by labels using name, version, and label filter.
//
// responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
//
// swagger:parameters removeByLabels
type RemoveByLabelsRequest struct {
	// Configuration group name
	// in: path
	// required: true
	GroupName string `json:"groupName"`

	// Configuration group version
	// in: path
	// required: true
	GroupVersion string `json:"groupVersion"`

	// Filter for configurations
	// in: path
	// required: true
	Filter string `json:"filter"`
}

func (c ConfigGroupHandler) RemoveByLabels(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["groupName"]
	groupVersion := mux.Vars(r)["groupVersion"]
	filterStr := mux.Vars(r)["filter"]

	groupVersionInt, err := strconv.Atoi(groupVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterMap := make(map[string]string)
	if filterStr != "" {
		keyValues := strings.Split(filterStr, ",")
		for _, kv := range keyValues {
			parts := strings.Split(kv, "=")
			if len(parts) == 2 {
				filterMap[parts[0]] = parts[1]
			}
		}
	}

	err = c.service.RemoveByLabels(groupName, groupVersionInt, filterMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:response ResponseConfigGroup
type ResponseConfigGroup struct {
	// Configuration group
	// in: body
	Body model.ConfigGroup `json:"body"`
}

// swagger:response ErrorResponse
type ErrorResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`

	// Message of the error
	// in: string
	Message string `json:"message"`
}

// swagger:response InternalServerErrorResponse
type InternalServerErrorResponse struct {
	// Error status code
	// in: int64
	Status int64 `json:"status"`

	// Message of the error
	// in: string
	Message string `json:"message"`
}
