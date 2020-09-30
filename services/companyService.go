package services

import (
	"Internship-Project/models"
)

var companies []models.Company = []models.Company{}

func AddNewCompany(newCompany models.Company) []models.Company {
	companies = append(companies, newCompany)

	return companies
}
