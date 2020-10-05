package controllers

import (
	"encoding/json"
	"internship_project/errorhandler"
	"internship_project/models"
	"internship_project/services"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetAllCompanies(w http.ResponseWriter, r *http.Request, connection *pgxpool.Pool) {
	companies, err := services.GetAllCompanies(connection)
	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func GetCompanyById(w http.ResponseWriter, r *http.Request, connection *pgxpool.Pool) {
	idParam := mux.Vars(r)["id"]

	company, err := services.GetCompany(idParam, connection)
	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func AddCompany(w http.ResponseWriter, r *http.Request, connection *pgxpool.Pool) {
	var newCompany models.Company
	json.NewDecoder(r.Body).Decode(&newCompany)
	err := services.AddNewCompany(&newCompany, connection)
	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCompany)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request, connection *pgxpool.Pool) {
	var idParam string = mux.Vars(r)["id"]

	var updateCompany models.Company
	json.NewDecoder(r.Body).Decode(&updateCompany)

	updateCompany.Id = idParam

	err := services.UpdateCompany(updateCompany, connection)

	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateCompany)

}

func DeleteCompany(w http.ResponseWriter, r *http.Request, connection *pgxpool.Pool) {
	var idParam string = mux.Vars(r)["id"]

	err := services.DeleteCompany(idParam, connection)

	if err != nil {
		writeErrToClient(w, err)
		return
	}
	w.WriteHeader(200)
}

func writeErrToClient(w http.ResponseWriter, err error) {
	errMsg, code := errorhandler.GetErrorMsg(err)
	w.WriteHeader(code)
	w.Write([]byte(errMsg))
}
