package repositories

import (
	"context"
	"errors"
	"fmt"
	"internship_project/models"

	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

//EmployeeRepository .
type EmployeeRepository struct {
	DB *pgxpool.Pool
}

// GetAllEmployees .
func (repository *EmployeeRepository) GetAllEmployees() ([]models.Employee, error) {
	allEmployees := []models.Employee{}
	rows, err := repository.DB.Query(context.Background(), "select * from employees")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var employee models.Employee
		err := rows.Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.CompanyID, &employee.C, &employee.R, &employee.U, &employee.D)
		if err != nil {
			return nil, err
		}
		allEmployees = append(allEmployees, employee)
	}
	return allEmployees, nil
}

// GetEmployeeByID .
func (repository *EmployeeRepository) GetEmployeeByID(id string) (models.Employee, error) {
	var employee models.Employee
	err := repository.DB.QueryRow(context.Background(), "select * from employees where id=$1", id).Scan(&employee.ID, &employee.FirstName, &employee.LastName, &employee.CompanyID, &employee.C, &employee.R, &employee.U, &employee.D)
	if err != nil {
		return employee, err
	}
	return employee, nil
}

// AddEmployee .
func (repository *EmployeeRepository) AddEmployee(employee *models.Employee) error {
	u := uuid.NewV4()
	employee.ID = u.String()
	_, err := repository.DB.Exec(context.Background(), "insert into employees (id, firstname, lastname, idc, c, r, u, d) values ($1, $2, $3, $4, $5, $6, $7, $8)",
		employee.ID, employee.FirstName, employee.LastName, employee.CompanyID, employee.C, employee.R, employee.U, employee.D)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEmployee .
func (repository *EmployeeRepository) UpdateEmployee(updatedEmp models.Employee) error {
	commandTag, err := repository.DB.Exec(context.Background(),
		"UPDATE employees SET firstname=$1, lastname=$2, idc=$3, c=$4, r=$5, u=$6, d=$7 WHERE id=$8;",
		updatedEmp.FirstName, updatedEmp.LastName, updatedEmp.CompanyID, updatedEmp.C, updatedEmp.R, updatedEmp.U, updatedEmp.D, updatedEmp.ID)

	if err != nil {
		fmt.Println(err)
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to update")
	}
	return nil
}

// DeleteEmployee .
func (repository *EmployeeRepository) DeleteEmployee(id string) error {
	commandTag, err := repository.DB.Exec(context.Background(), "DELETE FROM employees WHERE id=$1;", id)

	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("No row found to delete")
	}
	return nil
}
