package repositories

import (
	"context"
	"internship_project/models"
	"internship_project/persistence"
	"internship_project/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type CompanyRepository interface {
	GetAllCompanies() ([]models.Company, error)
	GetCompany(string) (models.Company, error)
	AddCompany(*models.Company) error
	UpdateCompany(models.Company) error
	DeleteCompany(string) error
	ChangeExternalRightApproveStatus(string, bool) error
}

type companyRepository struct {
	DB *pgxpool.Pool
}

func NewCompanyRepo(db *pgxpool.Pool) CompanyRepository {
	if db == nil {
		panic("CompanyRepository not created, pgxpool is nil")
	}
	return &companyRepository{
		DB: db,
	}
}

func (repository *companyRepository) GetAllCompanies() ([]models.Company, error) {
	companies := []models.Company{}
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
			ID:     stringUUID,
			Name:   company.Name,
			IsMain: company.Ismain,
		})
	}
	return companies, nil
}

func (repository *companyRepository) GetCompany(id string) (models.Company, error) {
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
		return company, utils.NoDataError
	}

	var companyPers persistence.Companies
	companyPers.Scan(&rows)

	var stringUUID string
	err = companyPers.Id.AssignTo(&stringUUID)
	if err != nil {
		return company, err
	}

	company = models.Company{
		ID:     stringUUID,
		Name:   companyPers.Name,
		IsMain: companyPers.Ismain,
	}

	return company, nil
}

func (repository *companyRepository) AddCompany(company *models.Company) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	company.ID = uuid.NewV4().String()
	companyPers := persistence.Companies{
		Name:   company.Name,
		Ismain: company.IsMain,
	}
	companyPers.Id.Set(company.ID)

	_, err = companyPers.InsertTx(&tx)
	if err != nil {
		return err
	}

	return tx.Commit(context.Background())
}

func (repository *companyRepository) UpdateCompany(company models.Company) error {
	tx, err := repository.DB.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	companyPers := persistence.Companies{
		Name:   company.Name,
		Ismain: company.IsMain,
	}
	companyPers.Id.Set(company.ID)

	commandTag, err := companyPers.UpdateTx(&tx)
	if err != nil {
		return err
	}
	if commandTag != 1 {
		return utils.NoDataError
	}

	return tx.Commit(context.Background())
}

func (repository *companyRepository) DeleteCompany(id string) error {
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
		return utils.NoDataError
	}
	return tx.Commit(context.Background())
}

func (repository *companyRepository) ChangeExternalRightApproveStatus(idear string, status bool) error {
	commandTag, err := repository.DB.Exec(context.Background(), "UPDATE external_access_rights SET approved = $1 WHERE id = $2;", status, idear)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return utils.NoDataError
	}
	return nil
}
