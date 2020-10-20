package controllers

import (
	"encoding/json"
	"internship_project/errorhandler"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
)

type EarController struct {
	Service services.EarService
}

func (controller *EarController) GetAllEars(w http.ResponseWriter, r *http.Request) {
	ears, err := controller.Service.GetAllEars()
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ears)
}

func (controller *EarController) GetEarById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	ear, err := controller.Service.GetEar(idParam)
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ear)
}

func (controller *EarController) AddEar(w http.ResponseWriter, r *http.Request) {
	var newEar models.ExternalRights
	json.NewDecoder(r.Body).Decode(&newEar)
	err := controller.Service.AddNewEar(&newEar)
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newEar)
}

func (controller *EarController) UpdateEar(w http.ResponseWriter, r *http.Request) {
	var updateEar models.ExternalRights
	json.NewDecoder(r.Body).Decode(&updateEar)

	err := controller.Service.UpdateEar(updateEar)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateEar)

}

func (controller *EarController) DeleteEar(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	err := controller.Service.DeleteEar(idParam)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(204)
}
