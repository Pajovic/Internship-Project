package controllers

import (
	"Internship-Project/models"
	"Internship-Project/services"
	"encoding/json"
	"net/http"
)

func AddCompany(w http.ResponseWriter, r *http.Request) {
	var newCompany models.Company
	json.NewDecoder(r.Body).Decode(&newCompany)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(services.AddNewCompany(newCompany))
}
