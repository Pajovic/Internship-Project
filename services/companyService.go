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

func GetCompany(id string, connection *pgx.Conn) (models.Company, error) {
	return repositories.GetCompany(id, connection)
}

func AddNewCompany(newCompany *models.Company, connection *pgx.Conn) error {
	return repositories.AddCompany(newCompany, connection)
}

func UpdateCompany(id int, updateCompany *models.Company) models.Company {
	//updateCompany.Id = id
	companies[id] = *updateCompany

	return companies[id]
}

func DeleteCompany(id int) {
	companies = append(companies[:id], companies[id+1:]...)
}
