package services

import (
	"internship_project/models"
	"internship_project/repositories"
)

//EmployeeService .
type EmployeeService struct {
	Repository repositories.EmployeeRepository
}

// GetAllEmployees is used to return all employees
func (service *EmployeeService) GetAllEmployees() ([]models.Employee, error) {
	return service.Repository.GetAllEmployees()
}

// AddNewEmployee is used to return all employees
func (service *EmployeeService) AddNewEmployee(newEmployee *models.Employee) error {
	return service.Repository.AddEmployee(newEmployee)
}

// GetEmployeeByID is used to find a specific employee
func (service *EmployeeService) GetEmployeeByID(id string) (models.Employee, error) {
	return service.Repository.GetEmployeeByID(id)
}

// UpdateEmployee is used to update a specific employee
func (service *EmployeeService) UpdateEmployee(updatedEmployee models.Employee) error {
	return service.Repository.UpdateEmployee(updatedEmployee)
}

// DeleteEmployee is used to update a specific employee
func (service *EmployeeService) DeleteEmployee(id string) error {
	return service.Repository.DeleteEmployee(id)
}
