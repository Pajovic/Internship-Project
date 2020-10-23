package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

type ExternalRightService struct {
	Repository repositories.ExternalRightRepository
}

func (service *ExternalRightService) GetAllEars() ([]models.ExternalRights, error) {
	return service.Repository.GetAllEars()
}

func (service *ExternalRightService) GetEar(id string) (models.ExternalRights, error) {
	return service.Repository.GetEar(id)
}

func (service *ExternalRightService) AddNewEar(newEar *models.ExternalRights) error {
	return service.Repository.AddEar(newEar)
}

func (service *ExternalRightService) UpdateEar(updateEar models.ExternalRights) error {
	return service.Repository.UpdateEar(updateEar)
}

func (service *ExternalRightService) DeleteEar(id string) error {
	return service.Repository.DeleteEar(id)
}
