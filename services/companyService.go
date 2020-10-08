package services

import (
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
