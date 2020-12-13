package services

import (
	"github.com/codingsince1985/geo-golang"
	"internship_project/models"
	"internship_project/repositories"
)

type ShopService struct {
	Repository repositories.ShopRepository
	Geocoder geo.Geocoder
}

func (service *ShopService) GetAllShops() ([]models.Shop, error) {
	return service.Repository.GetAllShops()
}

func (service *ShopService) SearchShopsByAddress(address string) ([]models.Shop, error) {
	location, err := service.Geocoder.Geocode(address)
	if err != nil {
		return nil, err
	}

	return service.Repository.GetShopsByLatLon(location.Lat, location.Lng)
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

func (service *ShopService) GetAddress(id string) (*geo.Address, error) {
	shop, err := service.Repository.GetShop(id)
	if err != nil {
		return nil, err
	}

	address, err := service.Geocoder.ReverseGeocode(shop.Lat, shop.Lon)
	if err != nil {
		return nil, err
	}

	return address, nil
}
