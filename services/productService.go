package services

import (
	"errors"
	"internship_project/models"
	"internship_project/repositories"
)

type ProductService struct {
	ProductRepository  repositories.ProductRepository
	EmployeeRepository repositories.EmployeeRepository
}

func (service *ProductService) GetAllProducts(employeeID string) ([]models.Product, error) {

	allProducts := []models.Product{}

	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return allProducts, err
	}

	if !employee.R {
		return allProducts, errors.New("You can't see products")
	}

	allProducts, err = service.ProductRepository.GetAllProducts(employee.CompanyID)

	if err != nil {
		return allProducts, err
	}

	return allProducts, nil
}

func (service *ProductService) GetProduct(productId string, employeeId string) (models.Product, error) {
	product := models.Product{}

	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeId)
	if err != nil {
		return product, err
	}

	if !employee.R {
		return product, errors.New("You can't see products")
	}

	product, err = service.ProductRepository.GetProduct(productId, employee.CompanyID)

	if err != nil {
		return product, err
	}

	return product, nil
}

func (service *ProductService) AddNewProduct(product *models.Product, employeeID string) error {
	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}

	if !employee.C {
		return errors.New("You can't create products")
	}

	if employee.CompanyID != product.IDC {
		return errors.New("You can't create products for other companies")
	}

	return service.ProductRepository.AddProduct(product)
}

func (service *ProductService) UpdateProduct(updateProduct models.Product, employeeId string) error {
	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeId)
	if err != nil {
		return err
	}

	if !employee.D {
		return errors.New("You can't update products")
	}

	product, err := service.ProductRepository.GetProduct(updateProduct.ID, employee.CompanyID)

	if employee.CompanyID != product.IDC {
		externalAccessRights, err := service.EmployeeRepository.GetEmployeeExternalPermissions(employee.CompanyID, product)
		if err != nil {
			return err
		}
		if !externalAccessRights.Update {
			return errors.New("You can't update this product")
		}
	}

	return service.ProductRepository.UpdateProduct(updateProduct)
}

func (service *ProductService) DeleteProduct(productId string, employeeId string) error {
	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeId)
	if err != nil {
		return err
	}

	if !employee.D {
		return errors.New("You can't delete products")
	}

	product, err := service.ProductRepository.GetProduct(productId, employee.CompanyID)

	if employee.CompanyID != product.IDC {
		externalAccessRights, err := service.EmployeeRepository.GetEmployeeExternalPermissions(employee.CompanyID, product)
		if err != nil {
			return err
		}
		if !externalAccessRights.Delete {
			return errors.New("You can't delete this product")
		}
	}

	return service.ProductRepository.DeleteProduct(productId)
}
