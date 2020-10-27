package repositories

import (
	"context"
	"internship_project/models"
	"internship_project/persistence"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type CompanyRepository struct {
	DB *pgxpool.Pool
}

func (repository *CompanyRepository) GetAllCompanies() ([]models.Company, error) {
	var companies []models.Company = []models.Company{}
	rows, err := repository.DB.Query(context.Background(), "select * from public.companies")
	defer rows.Close()

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
	return companies, nil
}

func (repository *CompanyRepository) GetCompany(id string) (models.Company, error) {
	var company models.Company

	Uuid, err := uuid.FromString(id)
	if err != nil {
		return company, err
	}

	rows, err := repository.DB.Query(context.Background(), `select * from public.companies where id = $1`, Uuid)
	defer rows.Close()

	if err != nil {
		return company, err
	}

	if !rows.Next() {
		var pgErr pgconn.PgError
		pgErr.Code = `02000`
		pgErr.Message = `There is no company with this id`
		return company, &pgErr
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

	return tx.Commit(context.Background())
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
		var pgErr pgconn.PgError
		pgErr.Code = `02000`
		pgErr.Message = `There is no company with this id`
		return &pgErr
	}

	return tx.Commit(context.Background())
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
		var pgErr pgconn.PgError
		pgErr.Code = `02000`
		pgErr.Message = `There is no company with this id`
		return &pgErr
	}
	return tx.Commit(context.Background())
}

func (repository *CompanyRepository) ChangeExternalRightApproveStatus(idear string, status bool) error {
	commandTag, err := repository.DB.Exec(context.Background(), "UPDATE external_access_rights SET approved = $1 WHERE id = $2;", status, idear)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		var pgErr pgconn.PgError
		pgErr.Code = `02000`
		pgErr.Message = `There is no external right with this id`
		return &pgErr
	}
	return nil
}
