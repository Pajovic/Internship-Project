package services

import (
	"internship_project/models"
	"internship_project/repositories"

	"github.com/jackc/pgx/v4/pgxpool"
)

var companies []models.Company = []models.Company{}

func GetAllCompanies(connection *pgxpool.Pool) ([]models.Company, error) {
	return repositories.GetAllCompanies(connection)
}

func GetCompany(id string, connection *pgxpool.Pool) (models.Company, error) {
	return repositories.GetCompany(id, connection)
}

func AddNewCompany(newCompany *models.Company, connection *pgxpool.Pool) error {
	return repositories.AddCompany(newCompany, connection)
}

func UpdateCompany(updateCompany models.Company, connection *pgxpool.Pool) error {
	return repositories.UpdateCompany(updateCompany, connection)
}

func DeleteCompany(id string, connection *pgxpool.Pool) error {
	return repositories.DeleteCompany(id, connection)
}
