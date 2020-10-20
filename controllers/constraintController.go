package controllers

import (
	"encoding/json"
	"internship_project/errorhandler"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
)

type ConstraintController struct {
	Service services.ConstraintService
}

func (controller *ConstraintController) GetAllConstraints(w http.ResponseWriter, r *http.Request) {
	constraints, err := controller.Service.GetAllConstraints()
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(constraints)
}

func (controller *ConstraintController) GetConstraintById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	constraint, err := controller.Service.GetConstraint(idParam)
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(constraint)
}

func (controller *ConstraintController) AddConstraint(w http.ResponseWriter, r *http.Request) {
	var newConstraint models.AccessConstraint
	json.NewDecoder(r.Body).Decode(&newConstraint)
	err := controller.Service.AddNewConstraint(&newConstraint)
	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newConstraint)
}

func (controller *ConstraintController) UpdateConstraint(w http.ResponseWriter, r *http.Request) {
	var updateConstraint models.AccessConstraint
	json.NewDecoder(r.Body).Decode(&updateConstraint)

	err := controller.Service.UpdateConstraint(updateConstraint)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateConstraint)

}

func (controller *ConstraintController) DeleteConstraint(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	err := controller.Service.DeleteConstraint(idParam)

	if err != nil {
		errorhandler.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(204)
}
