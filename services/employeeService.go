package services

import (
	"errors"
	"internship_project/models"
	"internship_project/repositories"
)

//EmployeeService .
type EmployeeService struct {
	Repository repositories.EmployeeRepository
}

// GetAllEmployees is used to return all employees
func (service *EmployeeService) GetAllEmployees(employeeID string) ([]models.Employee, error) {
	var allEmployees []models.Employee

	employee, err := service.Repository.GetEmployeeByID(employeeID)
	if err != nil {
		return allEmployees, err
	}

	if !employee.R {
		return allEmployees, errors.New("You have no permissions to preview other employees")
	}

	allEmployees, err = service.Repository.GetAllEmployees(employee.CompanyID)

	return allEmployees, nil
}

// AddNewEmployee is used to return all employees
func (service *EmployeeService) AddNewEmployee(newEmployee *models.Employee) error {
	return service.Repository.AddEmployee(newEmployee)
}

// GetEmployeeByID is used to find a specific employee
func (service *EmployeeService) GetEmployeeByID(id string, idEmployee string) (models.Employee, error) {
	employee, err := service.Repository.GetEmployeeByID(idEmployee)
	if err != nil {
		return models.Employee{}, err
	}

	employeeRequested, err := service.Repository.GetEmployeeByID(id)
	if err != nil {
		return models.Employee{}, err
	}

	if employee.CompanyID != employeeRequested.CompanyID {
		err := service.Repository.CheckCompaniesSharingEmployeeData(employee.CompanyID, employeeRequested.CompanyID)
		if err != nil {
			return models.Employee{}, err
		}
	}

	return employeeRequested, nil
}

// UpdateEmployee is used to update a specific employee
func (service *EmployeeService) UpdateEmployee(updatedEmployee models.Employee) error {
	return service.Repository.UpdateEmployee(updatedEmployee)
}

// DeleteEmployee is used to update a specific employee
func (service *EmployeeService) DeleteEmployee(id string) error {
	return service.Repository.DeleteEmployee(id)
}
