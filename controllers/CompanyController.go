package controllers

import (
	"Internship-Project/models"
	"Internship-Project/services"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
