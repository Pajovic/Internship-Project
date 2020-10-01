package repositories

import (
	"context"
	"internship_project/models"

	"github.com/jackc/pgx/v4"
	"github.com/satori/go.uuid"
)

func GetAllCompanies(connection *pgx.Conn) ([]models.Company, error) {
	var companies []models.Company = []models.Company{}
	rows, err := connection.Query(context.Background(), "select * from companies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var company models.Company
		err := rows.Scan(&company.Id, &company.Name, &company.IsMain)
		if err != nil {
			return nil, err
		}
		companies = append(companies, company)
	}
	return companies, nil
}

func GetCompany(id string, connection *pgx.Conn) (models.Company, error) {
	var company models.Company
	err := connection.QueryRow(context.Background(), "select * from companies where id=$1", id).Scan(&company.Id, &company.Name, &company.IsMain)
	if err != nil {
		return company, err
	}
	return company, nil
}

func AddCompany(company *models.Company, connection *pgx.Conn) error {
	u := uuid.NewV4()
	company.Id = u.String()
	_, err := connection.Exec(context.Background(), "insert into public.companies (id, name, ismain) values ($1, $2, $3)", u.Bytes(), company.Name, company.IsMain)
	if err != nil {
		return err
	}
//	if commandTag.RowsAffected() != 1 {
//		return errors.New("Error while inserting company")
//	}
	return nil
}
