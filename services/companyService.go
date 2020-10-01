package services

import (
	"internship_project/models"
	"internship_project/repositories"

	"github.com/jackc/pgx/v4"
)

var companies []models.Company = []models.Company{}

func GetAllCompanies(connection *pgx.Conn) ([]models.Company, error) {
	return repositories.GetAllCompanies(connection)
}

func GetCompany(id int) models.Company {
	return companies[id]
}

func AddNewCompany(newCompany *models.Company) []models.Company {
	//newCompany.Id = len(companies)
	companies = append(companies, *newCompany)

	return companies
}

func UpdateCompany(id int, updateCompany *models.Company) models.Company {
	//updateCompany.Id = id
	companies[id] = *updateCompany

	return companies[id]
}

func DeleteCompany(id int) {
	companies = append(companies[:id], companies[id+1:]...)
}
