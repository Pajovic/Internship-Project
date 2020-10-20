package repositories

import (
	"context"
	"errors"
	"internship_project/models"
	"internship_project/persistence"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type CompanyRepository struct {
	DB *pgxpool.Pool
}

func (repository *CompanyRepository) GetAllCompanies() ([]models.Company, error) {
	var companies []models.Company = []models.Company{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.companies")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var company persistence.Companies
		company.Scan(&rows)

		var stringUUID string
		err := company.Id.AssignTo(&stringUUID)
		if err != nil {
			return companies, err
		}

		companies = append(companies, models.Company{
			Id:     stringUUID,
			Name:   company.Name,
			IsMain: company.Ismain,
		})
	}
	rows.Close()
	return companies, nil
}

func (repository *CompanyRepository) GetCompany(id string) (models.Company, error) {
	var company models.Company

	Uuid, err := uuid.FromString(id)
	if err != nil {
		return company, err
	}

	rows, err := repository.DB.Query(context.Background(), `select * from companies where id = $1`, Uuid)
	defer rows.Close()

	if err != nil {
		return company, err
	}

	if !rows.Next() {
		return company, errors.New("There is no company with this id")
	}

	var companyPers persistence.Companies
	companyPers.Scan(&rows)

	var stringUUID string
	err = companyPers.Id.AssignTo(&stringUUID)
	if err != nil {
		return company, err
	}

	company = models.Company{
		Id:     stringUUID,
		Name:   companyPers.Name,
		IsMain: companyPers.Ismain,
	}

	return company, nil
}

func (repository *CompanyRepository) AddCompany(company *models.Company) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	company.Id = uuid.NewV4().String()
	companyPers := persistence.Companies{
		Name:   company.Name,
		Ismain: company.IsMain,
	}
	companyPers.Id.Set(company.Id)

	_, err = companyPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	tx.Commit(context.Background())

	return nil
}

func (repository *CompanyRepository) UpdateCompany(company models.Company) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	companyPers := persistence.Companies{
		Name:   company.Name,
		Ismain: company.IsMain,
	}
	companyPers.Id.Set(company.Id)

	commandTag, err := companyPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to update")
	}

	tx.Commit(context.Background())
	return nil
}

func (repository *CompanyRepository) DeleteCompany(id string) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	companyPers := persistence.Companies{}
	companyPers.Id.Set(id)

	commandTag, err := companyPers.DeleteTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return errors.New("No row found to delete")
	}
	tx.Commit(context.Background())
	return nil
}

func (repository *CompanyRepository) ApproveExternalAccess(idear string) error {
	commandTag, err := repository.DB.Exec(context.Background(), "UPDATE external_access_rights SET approved = true WHERE id = $1;", idear)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

func (repository *CompanyRepository) DisapproveExternalAccess(idear string) error {
	commandTag, err := repository.DB.Exec(context.Background(), "UPDATE external_access_rights SET approved = false WHERE id = $1;", idear)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}
