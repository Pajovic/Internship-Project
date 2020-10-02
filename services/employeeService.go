package services

import (
	"internship_project/models"
	"internship_project/repositories"

	"github.com/jackc/pgx/v4"
)

// GetAllEmployees is used to return all employees
func GetAllEmployees(conn *pgx.Conn) ([]models.Employee, error) {
	return repositories.GetAllEmployees(conn)
}

// AddNewEmployee is used to return all employees
func AddNewEmployee(conn *pgx.Conn, newEmployee *models.Employee) error {
	return repositories.AddEmployee(conn, newEmployee)
}

// GetEmployeeByID is used to find a specific employee
func GetEmployeeByID(conn *pgx.Conn, id string) (models.Employee, error) {
	return repositories.GetEmployeeByID(conn, id)
}

// UpdateEmployee is used to update a specific employee
func UpdateEmployee(conn *pgx.Conn, updatedEmployee models.Employee) error {
	return repositories.UpdateEmployee(conn, updatedEmployee)
}

// DeleteEmployee is used to update a specific employee
func DeleteEmployee(conn *pgx.Conn, id string) error {
	return repositories.DeleteEmployee(conn, id)
}
