package handlers

import (
	"encoding/json"
	"net/http"
	"projekat/model"
	"projekat/services"

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

// POST /config-groups/{name}/config
func (c ConfigGroupHandler) AddConfigToGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupName := vars["name"]

	var config model.Config
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = c.service.AddConfigToGroup(groupName, config) // Promenjeno c.service u c.groupService
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GET /configs
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
