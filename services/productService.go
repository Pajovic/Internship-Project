package services

import (
	"errors"
	"fmt"
	"internship_project/models"
	"internship_project/repositories"
)

type ProductService struct {
	ProductRepository  repositories.ProductRepository
	EmployeeRepository repositories.EmployeeRepository
}

func (service *ProductService) GetAllProducts(employeeID string) ([]models.Product, error) {
	var accessibleProducts []models.Product

	employee, err := service.EmployeeRepository.GetEmployeeByID(employeeID)
	if err != nil {
		fmt.Println("Error getting employee")
		return accessibleProducts, err
	}

	if !employee.R {
		return accessibleProducts, errors.New("You can't see products")
	}

	allProducts, err := service.ProductRepository.GetAllProducts()

	if err != nil {
		fmt.Println("Error getting products")
		return accessibleProducts, err
	}

	for _, product := range allProducts {
		if employee.CompanyID == product.IDC {
			// Product is within employee's company
			accessibleProducts = append(accessibleProducts, product)
		} else {
			// Product is owned by another company, we need to check access rights
			externalAccessRights, err := service.EmployeeRepository.GetEmployeeExternalPermissions(employee.CompanyID, product)
			if err != nil {
				//fmt.Println(err.Error() != "You don't have any permission for this product")
				if err.Error() != "You don't have any permission for this product" {
					return []models.Product{}, err
				}
				continue
			}
			if externalAccessRights.Read {
				accessibleProducts = append(accessibleProducts, product)
			}
		}
	}

	return accessibleProducts, nil
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
