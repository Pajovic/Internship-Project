package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

type ShopService struct {
	Repository repositories.ShopRepository
}

func (service *ShopService) GetAllShops() ([]models.Shop, error) {
	return service.Repository.GetAllShops()
}

func (service *ShopService) GetShop(id string) (models.Shop, error) {
	return service.Repository.GetShop(id)
}

func (service *ShopService) AddNewShop(newShop *models.Shop) error {
	return service.Repository.AddShop(newShop)
}

func (service *ShopService) UpdateShop(updateShop models.Shop) error {
	return service.Repository.UpdateShop(updateShop)
}

func (service *ShopService) DeleteShop(id string) error {
	return service.Repository.DeleteShop(id)
}
