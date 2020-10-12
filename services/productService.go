package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

type ProductService struct {
	Repository         repositories.ProductRepository
	EmployeeRepository repositories.EmployeeRepository
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

// GetEmployeePermissions is used to get employee's permissions towards a certain company
func (service *ProductService) GetEmployeePermissions(employeeID string, productID string) models.ExternalRights {
	employee, _ := service.EmployeeRepository.GetEmployeeByID(employeeID)

	/* if err_e != nil {
		return nil
	}
	if err_p != nil {
		return nil
	} */

	product, _ := service.GetProduct(productID)

	if employee.CompanyID == product.IDC {
		return models.ExternalRights{
			Read:   employee.R,
			Update: employee.U,
			Delete: employee.D,
		}
	}

	externalAccessRights, _ := service.EmployeeRepository.GetEmployeeExternalPermissions(employee.CompanyID, product)

	return externalAccessRights
}
