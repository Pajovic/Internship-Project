package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

type EarService struct {
	Repository repositories.EarRepository
}

func (service *EarService) GetAllEars() ([]models.ExternalRights, error) {
	return service.Repository.GetAllEars()
}

func (service *EarService) GetEar(id string) (models.ExternalRights, error) {
	return service.Repository.GetEar(id)
}

func (service *EarService) AddNewEar(newEar *models.ExternalRights) error {
	return service.Repository.AddEar(newEar)
}

func (service *EarService) UpdateEar(updateEar models.ExternalRights) error {
	return service.Repository.UpdateEar(updateEar)
}

func (service *EarService) DeleteEar(id string) error {
	return service.Repository.DeleteEar(id)
}
