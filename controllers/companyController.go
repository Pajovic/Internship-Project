package controllers

import (
	"encoding/json"
	"internship_project/models"
	"internship_project/services"
	"internship_project/utils"
	"net/http"

	"github.com/gorilla/mux"
)

type CompanyController struct {
	Service services.CompanyService
}

func (controller *CompanyController) GetAllCompanies(w http.ResponseWriter, r *http.Request) {
	companies, err := controller.Service.GetAllCompanies()
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func (controller *CompanyController) GetCompanyById(w http.ResponseWriter, r *http.Request) {
	idParam := mux.Vars(r)["id"]

	company, err := controller.Service.GetCompany(idParam)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

func (controller *CompanyController) AddCompany(w http.ResponseWriter, r *http.Request) {
	var newCompany models.Company
	json.NewDecoder(r.Body).Decode(&newCompany)
	err := controller.Service.AddNewCompany(&newCompany)
	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newCompany)
}

func (controller *CompanyController) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	var updateCompany models.Company
	json.NewDecoder(r.Body).Decode(&updateCompany)

	err := controller.Service.UpdateCompany(updateCompany)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updateCompany)
}

func (controller *CompanyController) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	var idParam string = mux.Vars(r)["id"]

	err := controller.Service.DeleteCompany(idParam)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(204)
}

func (controller *CompanyController) ChangeExternalRightApproveStatus(w http.ResponseWriter, r *http.Request, status bool) {
	var idear string = mux.Vars(r)["idear"]
	companyID := r.Header.Get("companyID")

	err := controller.Service.ChangeExternalRightApproveStatus(companyID, idear, status)

	if err != nil {
		utils.WriteErrToClient(w, err)
		return
	}
	w.WriteHeader(200)
}
