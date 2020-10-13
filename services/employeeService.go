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
func (service *EmployeeService) GetAllEmployees(employeeID string) ([]models.Employee, error) {
	var accessibleEmployees []models.Employee

	employee, err := service.Repository.GetEmployeeByID(employeeID)
	if err != nil {
		return accessibleEmployees, err
	}

	allEmployees, err := service.Repository.GetAllEmployees()

	for _, tempEmployee := range allEmployees {
		if tempEmployee.CompanyID == employee.CompanyID {
			// Employees are in the same company
			accessibleEmployees = append(accessibleEmployees, tempEmployee)
		} else {
			// Employees are in different companies, we need to check if they enabled sharing
			companiesSharingEmployeeData, err := service.Repository.CheckCompaniesSharingEmployeeData(employee.CompanyID, tempEmployee.CompanyID)
			if err != nil {
				if err.Error() != "You don't have any permission to view these employees" {
					return []models.Employee{}, err
				}
				continue
			}
			if companiesSharingEmployeeData {
				accessibleEmployees = append(accessibleEmployees, tempEmployee)
			}
		}
	}

	return accessibleEmployees, nil
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
