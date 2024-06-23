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
// Responses:
//
//	201: NoContent
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
func (c ConfigGroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var configGroup model.ConfigGroup
	err := json.NewDecoder(r.Body).Decode(&configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.Create(configGroup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// swagger:route GET /configGroups/{name}/{version} configGroups getConfigGroup
// Retrieves a configuration group by name and version.
//
// Responses:
//
//	200: ResponseConfigGroup
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
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
// Responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
func (c ConfigGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	version := mux.Vars(r)["version"]
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.Delete(name, versionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /configGroups configGroups getAllConfigGroups
// Retrieves all configuration groups.
//
// Responses:
//
//	200: []ResponseConfigGroup
//	500: InternalServerErrorResponse
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

// swagger:route DELETE /configGroups/{groupName}/{groupVersion}/removeConfig/{configName}/{configVersion} configGroups removeConfig
// Removes a configuration from a group by name and version.
//
// Responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	404: NotFoundResponse
//	500: InternalServerErrorResponse
func (c ConfigGroupHandler) RemoveConfig(w http.ResponseWriter, r *http.Request) {
	// Dohvatanje imena grupe, verzije grupe, imena konfiguracije i verzije konfiguracije iz putanje rute
	groupName := mux.Vars(r)["groupName"]
	groupVersion := mux.Vars(r)["groupVersion"]
	configName := mux.Vars(r)["configName"]
	configVersion := mux.Vars(r)["configVersion"]

	// Konverzija verzije grupe i verzije konfiguracije u integer
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

	// Poziv servisa za uklanjanje konfiguracije iz grupe
	err = c.service.RemoveConfig(groupName, groupVersionInt, configName, configVersionInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// swagger:route PUT /configGroups/{groupName}/{groupVersion}/addConfig configGroups addConfig
// Adds a configuration to a group by name and version.
//
// Responses:
//
//	201: NoContent
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
func (c ConfigGroupHandler) AddConfig(w http.ResponseWriter, r *http.Request) {
	// Dohvatanje imena grupe i verzije grupe iz putanje rute
	groupName := mux.Vars(r)["groupName"]
	groupVersion := mux.Vars(r)["groupVersion"]

	// Konverzija verzije grupe u integer
	groupVersionInt, err := strconv.Atoi(groupVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Dekodiranje tela zahteva kako bismo dobili objekat konfiguracije
	config := model.Config2{}
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Poziv servisa za dodavanje konfiguracije u grupu
	err = c.service.AddConfigs(groupName, groupVersionInt, config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFilteredConfigs handles GET requests to retrieve filtered configurations from a group.
//
// swagger:route GET /configGroups/{name}/{version}/filteredConfigs/{filter} configGroups getFilteredConfigs
// Retrieves filtered configurations from a group by name, version, and filter criteria.
//
// Responses:
//
//	200: []ResponseConfig
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
func (c ConfigGroupHandler) GetFilteredConfigs(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	versionStr := mux.Vars(r)["version"]
	filterStr := mux.Vars(r)["filter"]

	// Konvertovanje version iz stringa u int
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parsiranje filtera u mapu stringova
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

// swagger:route DELETE /configGroups/{groupName}/{groupVersion}/removeByLabels/{filter} configGroups removeByLabels
// Removes configurations from a group by labels using name, version, and label filter.
//
// Responses:
//
//	204: NoContent
//	400: BadRequestResponse
//	500: InternalServerErrorResponse
func (c ConfigGroupHandler) RemoveByLabels(w http.ResponseWriter, r *http.Request) {
	groupName := mux.Vars(r)["groupName"]
	groupVersion := mux.Vars(r)["groupVersion"]
	filterStr := mux.Vars(r)["filter"]

	// Konvertovanje verzije grupe u integer
	groupVersionInt, err := strconv.Atoi(groupVersion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parsiranje filtera u mapu stringova
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

	// swagger:route DELETE /configGroups/{groupName}/{groupVersion}/removeByLabels/{filter} configGroups removeByLabels
	// Removes configurations from a group by labels using name, version, and label filter.
	//
	// Responses:
	//   204: NoContent
	//   400: BadRequestResponse
	//   500: InternalServerErrorResponse
	err = c.service.RemoveByLabels(groupName, groupVersionInt, filterMap)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
