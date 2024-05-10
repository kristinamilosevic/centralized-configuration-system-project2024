package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/model"
	"projekat/services"
	"strconv"

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

// POST /configGroups
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

// GET /configGroups/{name}/{version}
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

// DELETE /configGroups/{name}/{version}
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

// GET /configGroups
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

// DELETE /configGroups/{groupName}/{groupVersion}/removeConfig/{configName}/{configVersion}
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

// PUT /configGroups/{groupName}/{groupVersion}/addConfig
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
