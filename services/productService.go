package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

type ProductService struct {
	Repository repositories.ProductRepository
}

func (service *ProductService) GetAllProducts() ([]models.Product, error) {
	return service.Repository.GetAllProducts()
}

func (service *ProductService) GetProduct(id string) (models.Product, error) {
	return service.Repository.GetProduct(id)
}

func (service *ProductService) AddNewProduct(product *models.Product) error {
	return service.Repository.AddProduct(product)
}

func (service *ProductService) Updateproduct(product models.Product) error {
	return service.Repository.UpdateProduct(product)
}

func (service *ProductService) DeleteProduct(id string) error {
	return service.Repository.DeleteProduct(id)
}
