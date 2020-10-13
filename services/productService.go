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

func (service *ProductService) GetAllProducts() ([]models.Product, error) {
	return service.ProductRepository.GetAllProducts()
}

func (service *ProductService) GetProduct(productID string, employeeID string) (models.Product, error) {
	product, err := service.ProductRepository.GetProduct(productID)
	if err != nil {
		return product, err
	}

	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return product, err
	}

	if !employee.R {
		return product, errors.New("You can't see products")
	}

	if employee.CompanyID != product.IDC {
		externalAccessRights, err := service.EmployeeRepository.GetEmployeeExternalPermissions(employee.CompanyID, product)
		if err != nil {
			return product, err
		}
		if !externalAccessRights.Read {
			return product, errors.New("You can't see this product")
		}
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

func (service *ProductService) Updateproduct(updateProduct models.Product, employeeID string) error {
	product, err := service.ProductRepository.GetProduct(updateProduct.ID)
	if err != nil {
		return err
	}

	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}

	if !employee.U {
		return errors.New("You can't update products")
	}

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

func (service *ProductService) DeleteProduct(productID string, employeeID string) error {
	product, err := service.ProductRepository.GetProduct(productID)
	if err != nil {
		return err
	}

	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		return err
	}

	if !employee.D {
		return errors.New("You can't delete products")
	}

	if employee.CompanyID != product.IDC {
		externalAccessRights, err := service.EmployeeRepository.GetEmployeeExternalPermissions(employee.CompanyID, product)
		if err != nil {
			return err
		}
		if !externalAccessRights.Delete {
			return errors.New("You can't delete this product")
		}
	}

	return service.ProductRepository.DeleteProduct(productID)
}
