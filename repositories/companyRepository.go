package repositories

import (
	"context"
	"errors"
	"internship_project/models"

	"github.com/jackc/pgx/v4"
)

func GetAllCompanies(connection *pgx.Conn) ([]models.Company, error) {
	var companies []models.Company = []models.Company{}
	rows, err := connection.Query(context.Background(), "select * from companies")
	if err != nil {
		return nil, errors.New("Failed to get companies from database")
	}
	defer rows.Close()
	for rows.Next() {
		var company models.Company
		err := rows.Scan(&company.Id, &company.Name, &company.IsMain)
		if err != nil {
			return nil, errors.New("Failed to parse company from database")
		}
		companies = append(companies, company)
	}
	return companies, nil
}
