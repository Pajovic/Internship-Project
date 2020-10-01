package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func GetAllCompanies(w http.ResponseWriter, r *http.Request, connection *pgx.Conn) {
	companies, err := services.GetAllCompanies(connection)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error while getting companies"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func GetCompanyById(w http.ResponseWriter, r *http.Request, connection *pgx.Conn) {
	idParam := mux.Vars(r)["id"]

	company, err := services.GetCompany(idParam, connection)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("There is no company with this id"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func AddCompany(w http.ResponseWriter, r *http.Request, connection *pgx.Conn) {
	var newCompany models.Company
	json.NewDecoder(r.Body).Decode(&newCompany)
	err := services.AddNewCompany(&newCompany, connection)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte("Error while inserting company"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCompany)
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
