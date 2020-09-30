package services

import (
	"internship_project/models"
)

var companies []models.Company = []models.Company{}

func GetAllCompanies() []models.Company {
	return companies
}

func GetCompany(id int) models.Company {
	return companies[id]
}

func AddNewCompany(newCompany *models.Company) []models.Company {
	newCompany.Id = len(companies)
	companies = append(companies, *newCompany)

	return companies
}

func UpdateCompany(id int, updateCompany *models.Company) models.Company {
	updateCompany.Id = id
	companies[id] = *updateCompany

	return companies[id]
}

func DeleteCompany(id int) {
	companies = append(companies[:id], companies[id+1:]...)
}
