package services

import (
	"errors"
	"internship_project/models"
	"internship_project/repositories"
)

type CompanyService struct {
	Repository repositories.CompanyRepository
}

func (service *CompanyService) GetAllCompanies() ([]models.Company, error) {
	return service.Repository.GetAllCompanies()
}

func (service *CompanyService) GetCompany(id string) (models.Company, error) {
	return service.Repository.GetCompany(id)
}

func (service *CompanyService) AddNewCompany(newCompany *models.Company) error {
	return service.Repository.AddCompany(newCompany)
}

func (service *CompanyService) UpdateCompany(updateCompany models.Company) error {
	return service.Repository.UpdateCompany(updateCompany)
}

func (service *CompanyService) DeleteCompany(id string) error {
	return service.Repository.DeleteCompany(id)
}

func (service *CompanyService) ApproveExternalAccess(companyID string, idear string) error {
	approvingCompany, err := service.Repository.GetCompany(companyID)
	if err != nil {
		return err
	}

	if !approvingCompany.IsMain {
		return errors.New("Your company does not have permission to approve sharing")
	}

	return service.Repository.ApproveExternalAccess(idear)
}

func (service *CompanyService) DisapproveExternalAccess(companyID string, idear string) error {
	approvingCompany, err := service.Repository.GetCompany(companyID)
	if err != nil {
		return err
	}

	if !approvingCompany.IsMain {
		return errors.New("Your company does not have permission to disapprove sharing")
	}

	return service.Repository.DisapproveExternalAccess(idear)
}
