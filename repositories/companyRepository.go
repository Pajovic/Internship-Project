package repositories

import (
	"context"
	"errors"
	"fmt"
	"internship_project/models"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

func GetAllCompanies(connection *pgxpool.Pool) ([]models.Company, error) {
	var companies []models.Company = []models.Company{}
	rows, err := connection.Query(context.Background(), "select * from companies")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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

func GetCompany(id string, connection *pgxpool.Pool) (models.Company, error) {
	var company models.Company
	err := connection.QueryRow(context.Background(), "select * from companies where id=$1", id).Scan(&company.Id, &company.Name, &company.IsMain)
	if err != nil {
		return company, err
	}
	return company, nil
}

func AddCompany(company *models.Company, connection *pgxpool.Pool) error {
	u := uuid.NewV4()
	company.Id = u.String()
	_, err := connection.Exec(context.Background(), "insert into public.companies (id, name, ismain) values ($1, $2, $3)", u.Bytes(), company.Name, company.IsMain)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCompany(company models.Company, conn *pgxpool.Pool) error {
	commandTag, err := conn.Exec(context.Background(),
		"UPDATE public.companies SET name=$1, ismain=$2 WHERE id=$3",
		company.Name, company.IsMain, company.Id)

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

func DeleteCompany(id string, conn *pgxpool.Pool) error {
	commandTag, err := conn.Exec(context.Background(), "DELETE FROM public.companies WHERE id=$1;", id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}
