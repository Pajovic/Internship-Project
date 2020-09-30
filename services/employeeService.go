package services

import "internship_project/models"

// GetAllEmployees is used to return all employees
func GetAllEmployees(employees []models.Employee) []models.Employee {
	return employees
}

// AddNewEmployee is used to return all employees
func AddNewEmployee(employees *[]models.Employee, newEmployee models.Employee) (*[]models.Employee, error) {
	*employees = append(*employees, newEmployee)
	return employees, nil
}

// GetEmployeeByID is used to find a specific employee
func GetEmployeeByID(employees []models.Employee, id string) models.Employee {
	employee := models.Employee{}

	for i := range employees {
		if employees[i].ID == id {
			employee = employees[i]
			break
		}
	}

	return employee
}

// UpdateEmployee is used to update a specific employee
func UpdateEmployee(employees *[]models.Employee, updatedEmployee models.Employee, id string) []models.Employee {
	for i := range *employees {
		if (*employees)[i].ID == id {
			(*employees)[i] = updatedEmployee
			break
		}
	}

	return *employees
}

// DeleteEmployee is used to update a specific employee
func DeleteEmployee(employees []models.Employee, id string) []models.Employee {
	var index int

	for i := range employees {
		if employees[i].ID == id {
			index = i
			break
		}
	}

	return append(employees[:index], employees[index+1:]...)
}
