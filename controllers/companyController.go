package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.GetAllCompanies())
}

func GetCompanyById(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.GetCompany(id))
}

func AddCompany(w http.ResponseWriter, r *http.Request) {
	var newCompany models.Company
	json.NewDecoder(r.Body).Decode(&newCompany)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.AddNewCompany(&newCompany))
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	var updateCompany models.Company
	json.NewDecoder(r.Body).Decode(&updateCompany)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.UpdateCompany(id, &updateCompany))

}

func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("ID could not be converted to integer"))
		return
	}

	services.DeleteCompany(id)
	w.WriteHeader(200)
}
