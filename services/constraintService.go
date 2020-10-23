package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

type ConstraintService struct {
	Repository repositories.ConstraintRepository
}

func (service *ConstraintService) GetAllConstraints() ([]models.AccessConstraint, error) {
	return service.Repository.GetAllConstraints()
}

func (service *ConstraintService) GetConstraint(id string) (models.AccessConstraint, error) {
	return service.Repository.GetConstraint(id)
}

func (service *ConstraintService) AddNewConstraint(newConstraint *models.AccessConstraint) error {
	return service.Repository.AddConstraint(newConstraint)
}

func (service *ConstraintService) UpdateConstraint(updateConstraint models.AccessConstraint) error {
	return service.Repository.UpdateConstraint(updateConstraint)
}

func (service *ConstraintService) DeleteConstraint(id string) error {
	return service.Repository.DeleteConstraint(id)
}
